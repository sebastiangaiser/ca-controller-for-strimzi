#!/usr/bin/env bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

NAMESPACE="${NAMESPACE:-kafka}"
TEST_PASSED=0
TEST_FAILED=0

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

assert_equals() {
    local expected="$1"
    local actual="$2"
    local message="$3"

    if [[ "$expected" == "$actual" ]]; then
        log_info "PASS: $message"
        ((TEST_PASSED++)) || true
    else
        log_error "FAIL: $message (expected: '$expected', actual: '$actual')"
        ((TEST_FAILED++)) || true
    fi
}

assert_not_empty() {
    local value="$1"
    local message="$2"

    if [[ -n "$value" ]]; then
        log_info "PASS: $message"
        ((TEST_PASSED++)) || true
    else
        log_error "FAIL: $message (value is empty)"
        ((TEST_FAILED++)) || true
    fi
}

assert_secret_exists() {
    local secret_name="$1"
    local namespace="$2"

    if kubectl get secret "$secret_name" -n "$namespace" &>/dev/null; then
        log_info "PASS: Secret '$secret_name' exists in namespace '$namespace'"
        ((TEST_PASSED++)) || true
    else
        log_error "FAIL: Secret '$secret_name' does not exist in namespace '$namespace'"
        ((TEST_FAILED++)) || true
    fi
}

assert_secret_not_exists() {
    local secret_name="$1"
    local namespace="$2"

    if ! kubectl get secret "$secret_name" -n "$namespace" &>/dev/null; then
        log_info "PASS: Secret '$secret_name' does not exist in namespace '$namespace' (as expected)"
        ((TEST_PASSED++)) || true
    else
        log_error "FAIL: Secret '$secret_name' should not exist in namespace '$namespace'"
        ((TEST_FAILED++)) || true
    fi
}

wait_for_secret() {
    local secret_name="$1"
    local namespace="$2"
    local timeout="${3:-60}"

    log_info "Waiting for secret '$secret_name' in namespace '$namespace' (timeout: ${timeout}s)..."

    local count=0
    while [[ $count -lt $timeout ]]; do
        if kubectl get secret "$secret_name" -n "$namespace" &>/dev/null; then
            log_info "Secret '$secret_name' is ready"
            return 0
        fi
        sleep 1
        ((count++))
    done

    log_error "Timeout waiting for secret '$secret_name'"
    return 1
}

cleanup_test_resources() {
    log_info "Cleaning up test resources..."
    kubectl delete certificate my-cluster-clients-ca-cert-tls -n "$NAMESPACE" --ignore-not-found
    kubectl delete certificate my-cluster-cluster-ca-cert-tls -n "$NAMESPACE" --ignore-not-found
    kubectl delete secret my-cluster-clients-ca-cert -n "$NAMESPACE" --ignore-not-found
    kubectl delete secret my-cluster-clients-ca -n "$NAMESPACE" --ignore-not-found
    kubectl delete secret my-cluster-cluster-ca-cert -n "$NAMESPACE" --ignore-not-found
    kubectl delete secret my-cluster-cluster-ca -n "$NAMESPACE" --ignore-not-found
    # Clean up any historical secrets
    kubectl delete secrets -n "$NAMESPACE" -l sebastian.gaiser.bayern/historical=true --ignore-not-found
    sleep 2
}

# ==============================================================================
# Test: Controller is running
# ==============================================================================
test_controller_running() {
    log_info "=== Test: Controller is running ==="

    local ready
    ready=$(kubectl get deployment ca-controller-for-strimzi -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
    assert_equals "1" "$ready" "Controller deployment has 1 ready replica"
}

# ==============================================================================
# Test: Create certificates and verify target secrets
# ==============================================================================
test_certificate_creation() {
    log_info "=== Test: Certificate creation ==="

    # Apply the ClusterIssuer first
    log_info "Creating ClusterIssuer..."
    cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-cluster-issuer
spec:
  selfSigned: {}
EOF

    # Create a root CA
    log_info "Creating root CA certificate..."
    cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-ca
  namespace: cert-manager
spec:
  isCA: true
  commonName: example-ca
  secretName: root-secret
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: selfsigned-cluster-issuer
    kind: ClusterIssuer
    group: cert-manager.io
EOF

    # Wait for root CA to be ready
    kubectl wait --for=condition=Ready certificate/example-ca -n cert-manager --timeout=60s

    # Create the CA issuer
    log_info "Creating CA ClusterIssuer..."
    cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: my-ca
spec:
  ca:
    secretName: root-secret
EOF

    sleep 5 # Wait for issuer to be ready

    # Create the test certificate
    log_info "Creating test certificate..."
    cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-cluster-clients-ca-cert-tls
  namespace: $NAMESPACE
spec:
  isCA: true
  commonName: my-cluster-clients-ca-cert-tls
  secretName: my-cluster-clients-ca-cert-tls
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: my-ca
    kind: ClusterIssuer
    group: cert-manager.io
  secretTemplate:
    annotations:
      sebastian.gaiser.bayern/tls-strimzi-ca: "reconcile"
      sebastian.gaiser.bayern/target-cluster-name: "my-cluster"
      sebastian.gaiser.bayern/target-secret-name: "my-cluster-clients-ca-cert"
      sebastian.gaiser.bayern/target-secret-key-name: "my-cluster-clients-ca"
EOF

    # Wait for certificate to be ready
    kubectl wait --for=condition=Ready certificate/my-cluster-clients-ca-cert-tls -n "$NAMESPACE" --timeout=60s

    # Wait for controller to create target secrets
    wait_for_secret "my-cluster-clients-ca-cert" "$NAMESPACE" 30
    wait_for_secret "my-cluster-clients-ca" "$NAMESPACE" 30

    # Verify target secrets exist
    assert_secret_exists "my-cluster-clients-ca-cert" "$NAMESPACE"
    assert_secret_exists "my-cluster-clients-ca" "$NAMESPACE"

    # Verify secret labels
    local managed_by
    managed_by=$(kubectl get secret my-cluster-clients-ca-cert -n "$NAMESPACE" -o jsonpath='{.metadata.labels.sebastian\.gaiser\.bayern/managed-by}')
    assert_equals "ca-controller-for-strimzi" "$managed_by" "Target secret has correct managed-by label"

    # Verify secret has ca.crt data
    local ca_crt
    ca_crt=$(kubectl get secret my-cluster-clients-ca-cert -n "$NAMESPACE" -o jsonpath='{.data.ca\.crt}')
    assert_not_empty "$ca_crt" "Target secret has ca.crt data"

    # Verify key secret has ca.key data
    local ca_key
    ca_key=$(kubectl get secret my-cluster-clients-ca -n "$NAMESPACE" -o jsonpath='{.data.ca\.key}')
    assert_not_empty "$ca_key" "Key secret has ca.key data"

    # Verify generation annotation
    local generation
    generation=$(kubectl get secret my-cluster-clients-ca-cert -n "$NAMESPACE" -o jsonpath='{.metadata.annotations.strimzi\.io/ca-cert-generation}')
    assert_equals "0" "$generation" "Initial generation is 0"
}

# ==============================================================================
# Test: Certificate rotation creates historical secrets
# ==============================================================================
test_certificate_rotation() {
    log_info "=== Test: Certificate rotation and historical secrets ==="

    # Get current hash before rotation
    local hash_before
    hash_before=$(kubectl get secret my-cluster-clients-ca-cert -n "$NAMESPACE" -o jsonpath='{.metadata.labels.sebastian\.gaiser\.bayern/hash}')
    log_info "Hash before rotation: $hash_before"

    # Trigger certificate renewal by deleting the source secret
    # cert-manager will recreate it with new cert data
    log_info "Triggering certificate rotation..."
    kubectl delete secret my-cluster-clients-ca-cert-tls -n "$NAMESPACE"

    # Wait for cert-manager to recreate the secret
    sleep 5
    kubectl wait --for=condition=Ready certificate/my-cluster-clients-ca-cert-tls -n "$NAMESPACE" --timeout=60s

    # Wait for controller to process the change
    sleep 10

    # Check for historical secrets
    log_info "Checking for historical secrets..."
    local historical_cert_secret="my-cluster-clients-ca-cert-generation-0"
    local historical_key_secret="my-cluster-clients-ca-generation-0"

    # Wait a bit more for the controller to create historical secrets
    sleep 5

    # Verify historical secrets were created
    assert_secret_exists "$historical_cert_secret" "$NAMESPACE"
    assert_secret_exists "$historical_key_secret" "$NAMESPACE"

    # Verify historical secret has the historical label
    local is_historical
    is_historical=$(kubectl get secret "$historical_cert_secret" -n "$NAMESPACE" -o jsonpath='{.metadata.labels.sebastian\.gaiser\.bayern/historical}' 2>/dev/null || echo "")
    assert_equals "true" "$is_historical" "Historical secret has historical=true label"

    # Verify the main secret was updated with new generation
    local new_generation
    new_generation=$(kubectl get secret my-cluster-clients-ca-cert -n "$NAMESPACE" -o jsonpath='{.metadata.annotations.strimzi\.io/ca-cert-generation}')
    assert_equals "1" "$new_generation" "Generation incremented to 1 after rotation"

    # Verify hash changed
    local hash_after
    hash_after=$(kubectl get secret my-cluster-clients-ca-cert -n "$NAMESPACE" -o jsonpath='{.metadata.labels.sebastian\.gaiser\.bayern/hash}')
    log_info "Hash after rotation: $hash_after"

    if [[ "$hash_before" != "$hash_after" ]]; then
        log_info "PASS: Hash changed after rotation"
        ((TEST_PASSED++)) || true
    else
        log_error "FAIL: Hash did not change after rotation"
        ((TEST_FAILED++)) || true
    fi
}

# ==============================================================================
# Main
# ==============================================================================
main() {
    log_info "Starting e2e tests..."
    log_info "Namespace: $NAMESPACE"

    # Cleanup before tests
    cleanup_test_resources

    # Run tests
    test_controller_running
    test_certificate_creation
    test_certificate_rotation

    # Cleanup after tests
    cleanup_test_resources

    # Print summary
    echo ""
    log_info "======================================"
    log_info "Test Summary"
    log_info "======================================"
    log_info "Passed: $TEST_PASSED"
    if [[ $TEST_FAILED -gt 0 ]]; then
        log_error "Failed: $TEST_FAILED"
        exit 1
    else
        log_info "Failed: $TEST_FAILED"
        log_info "All tests passed!"
        exit 0
    fi
}

main "$@"
