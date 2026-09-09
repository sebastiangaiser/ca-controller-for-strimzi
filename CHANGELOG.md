# Changelog

## [v0.5.2](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.5.2)

[Compare to previous version](https://github.com/sebastiangaiser/ca-controller-for-strimzi/compare/v0.5.1...v0.5.2)

### Bug Fixes

- **deps**: update kubernetes monorepo to v0.35.4 (#197) ([3e45903](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/3e45903f775ad4b4baa64cef6d33f8b8538c7fa1))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.28.2 (#200) ([931d2bd](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/931d2bd801442b6058f49baaed9e3a07390e78fe))
- **deps**: update module go.uber.org/zap to v1.28.0 (#201) ([c357a54](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/c357a545997e7b71dfc58d6f35ded0704d29b436))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.28.3 (#202) ([58e2cda](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/58e2cda79ea6b4772ce6042622e5c2f7c33b6bcc))
- **deps**: update module sigs.k8s.io/controller-runtime to v0.24.0 (#203) ([b4bd9d1](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/b4bd9d14674f887e2522663b911e326246a9068e))
- **deps**: update module sigs.k8s.io/controller-runtime to v0.24.1 (#204) ([c24e523](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/c24e5232b15b605ae10eb2e48de500c72cd3fbe2))
- **deps**: update kubernetes monorepo to v0.36.1 (#205) ([31a0a68](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/31a0a6886b896e129ad93a39855e699570878381))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.29.0 (#206) ([4dd2270](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/4dd22701eab8324fbb5aaeaccc0aad21cf2e33ea))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.30.0 (#210) ([2c7d0f7](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/2c7d0f737cc4e3e974f30edee946f5efe66a71d6))
- **deps**: update kubernetes monorepo to v0.36.2 (#208) ([66fc02f](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/66fc02f7c070e8a111560d688b220aaa4e4f934a))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.31.0 (#211) ([34070e9](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/34070e986ed3bfac4f5df141f6d90745f5930af4))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.32.0 (#214) ([e2a371a](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/e2a371a4040dfff5a9289aee58b5d52c7a385036))
- **deps**: update module github.com/prometheus/client_golang to v1.24.0 (#222) ([7e3e788](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/7e3e78871db98cd0c5b6f5ad72de3801cc89174e))
- **deps**: update kubernetes monorepo to v0.36.3 (#223) ([1aaab7f](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/1aaab7f20e523c396c485882a55d72e75fbe41ce))
- **deps**: update module github.com/prometheus/client_golang to v1.24.1 (#224) ([b858b8e](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/b858b8ef226641d01830a610b6d1aee6e00f248c))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.32.1 (#226) ([f0571db](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/f0571dbdefad0a6fb9c42170151e749b2382bbc5))
- **deps**: update kubernetes monorepo to v0.36.4 (#230) ([ff76970](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/ff769703ae2e1ac2d894ad6b20262af2342e3c06))
- **deps**: update kubernetes monorepo to v0.37.0 (#231) ([a19e9ce](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/a19e9ce685e8e1acd1a0728a75e22d77dbc1f401))
- **deps**: update module sigs.k8s.io/controller-runtime to v0.25.0 (#233) ([b169607](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/b16960775559de00fb1f0b846d7cbab8ea7bd99e))
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.32.2 (#234) ([3887203](https://github.com/sebastiangaiser/ca-controller-for-strimzi/commit/38872030c078a5b3a220c879e78be75c977ae374))

## [v0.5.1](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.5.1)

### Bug Fixes

- restore combined ca.crt (cluster CA + issuing CA) for Strimzi signing compatibility (#192)

## [v0.5.0](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.5.0)

### Features

- support rotationPolicy=Never for not rotating the private key (#189)

### Bug Fixes

- **deps**: update kubernetes monorepo to v0.35.3 (#190)

## [v0.4.0](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.4.0)

### Features

- add secret history (#174)

### Bug Fixes

- **deps**: update module github.com/onsi/ginkgo/v2 to v2.27.1 (#150)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.27.2 (#152)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.22.4 (#153)
- **deps**: update kubernetes packages to v0.34.2 (#155)
- **deps**: update module go.uber.org/zap to v1.27.1 (#156)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.27.3 (#160)
- **deps**: update kubernetes packages to v0.34.3 (#161)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.27.4 (#164)
- **deps**: update kubernetes packages to v0.35.0 (#162)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.27.5 (#165)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.23.0 (#167)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.23.1 (#168)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.28.1 (#170)
- **deps**: update kubernetes packages to v0.35.1 (#172)
- **deps**: update kubernetes packages to v0.35.2 (#185)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.23.3 (#187)
- use 'ca.crt' and 'tls.crt' directly instead of combining them (#188)

## [v0.3.0](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.3.0)

### Features

- **chart**: add CiliumNetworkPolicy (#146)

## [v0.2.1](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.2.1)

### Bug Fixes

- **deps**: update kubernetes packages to v0.33.2 (#108)
- **deps**: update kubernetes packages to v0.33.3 (#114)
- **deps**: update module github.com/prometheus/client_golang to v1.23.0 (#117)
- **deps**: update kubernetes packages to v0.33.4 (#122)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.24.0 (#123)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.25.0 (#125)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.25.1 (#126)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.22.0 (#129)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.25.2 (#130)
- **deps**: update module github.com/prometheus/client_golang to v1.23.2 (#134)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.25.3 (#133)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.22.1 (#136)
- **deps**: update kubernetes packages to v0.34.1 (#137)
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.26.0 (#139)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.22.2 (#140)
- **deps**: update module sigs.k8s.io/controller-runtime to v0.22.3 (#144)

## [v0.2.0](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.2.0)

### Features

- **helm**: set default securityContext matching pod-security &#39;restricted&#39; (#103)

## [v0.1.11](https://github.com/sebastiangaiser/ca-controller-for-strimzi/releases/tag/v0.1.11)

### Bug Fixes

- **deps**: update module sigs.k8s.io/controller-runtime to v0.20.1
- **deps**: update kubernetes packages to v0.32.2
- **deps**: update module sigs.k8s.io/controller-runtime to v0.20.2
- **deps**: update module github.com/prometheus/client_golang to v1.21.0
- **deps**: update module github.com/prometheus/client_golang to v1.21.1
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.23.0
- **deps**: update module sigs.k8s.io/controller-runtime to v0.20.3
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.23.1
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.23.2
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.23.3
- **deps**: update module sigs.k8s.io/controller-runtime to v0.20.4
- **deps**: update kubernetes packages to v0.32.3
- **deps**: update module github.com/onsi/ginkgo/v2 to v2.23.4
- **deps**: update module github.com/prometheus/client_golang to v1.22.0
- **deps**: update kubernetes packages to v0.32.4
- **deps**: update kubernetes packages to v0.33.0
- **deps**: update kubernetes packages to v0.33.1
- **deps**: update module sigs.k8s.io/controller-runtime to v0.21.0
