# DevContainer Multi-Architecture

Ce DevContainer supporte automatiquement les architectures **AMD64 (x86_64)** et **ARM64 (aarch64)**.

## Architecture Automatique

Le Dockerfile détecte automatiquement l'architecture de la plateforme hôte grâce à la variable Docker `TARGETARCH` :

```dockerfile
ARG TARGETARCH
```

### Composants avec détection d'architecture

1. **AWS CLI v2**
   - ARM64 : `awscli-exe-linux-aarch64.zip`
   - AMD64 : `awscli-exe-linux-x86_64.zip`

2. **Go 1.25.2**
   - ARM64 : `go1.25.2.linux-arm64.tar.gz`
   - AMD64 : `go1.25.2.linux-amd64.tar.gz`

3. **Packages APT**
   - Détection automatique via `$(dpkg --print-architecture)`

## Utilisation Locale

### VS Code DevContainer
Le DevContainer s'adapte automatiquement à votre architecture :
- Mac M1/M2/M3 → ARM64
- Mac Intel / PC → AMD64

### Docker Build Manuel
```bash
# Build pour l'architecture de votre machine
docker build -t ktn-linter-dev .devcontainer/

# Build pour une architecture spécifique
docker build --platform linux/amd64 -t ktn-linter-dev:amd64 .devcontainer/
docker build --platform linux/arm64 -t ktn-linter-dev:arm64 .devcontainer/
```

## Utilisation en CI/CD

### GitHub Actions
Le workflow `.github/workflows/build-devcontainer.yml` build automatiquement pour les deux architectures :

```yaml
strategy:
  matrix:
    platform:
      - linux/amd64
      - linux/arm64
```

### Build Multi-Architecture avec Docker Buildx

```bash
# Créer un builder multi-architecture
docker buildx create --use --name multiarch

# Build et push pour plusieurs architectures
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/your-org/ktn-linter:latest \
  --push \
  .devcontainer/
```

## Pipeline GitLab CI

```yaml
build-devcontainer:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker buildx create --use
  script:
    - docker buildx build
        --platform linux/amd64,linux/arm64
        -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
        --push
        .devcontainer/
  parallel:
    matrix:
      - ARCH: [amd64, arm64]
```

## AWS CodeBuild

```yaml
version: 0.2

phases:
  pre_build:
    commands:
      - echo Logging in to Amazon ECR...
      - aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com
      - docker buildx create --use
  build:
    commands:
      - echo Building multi-arch Docker image...
      - docker buildx build --platform linux/amd64,linux/arm64 -t $IMAGE_REPO_NAME:$IMAGE_TAG --push .devcontainer/
```

## Vérification de l'Architecture

Pour vérifier quelle architecture a été construite :

```bash
# Dans le container
uname -m
# ARM64 : aarch64
# AMD64 : x86_64

# Version d'AWS CLI
aws --version

# Version de Go
go version
```

## Dépannage

### Erreur : "exec format error"
L'architecture du container ne correspond pas à votre système. Utilisez `--platform` :
```bash
docker run --platform linux/arm64 ktn-linter-dev
```

### Build QEMU lent
Les builds cross-architecture utilisent l'émulation QEMU qui peut être lente. Pour de meilleures performances, buildez nativement sur l'architecture cible.

### Cache Docker Buildx
Pour optimiser les builds multi-architecture :
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --cache-from type=registry,ref=myregistry.io/myapp:cache \
  --cache-to type=registry,ref=myregistry.io/myapp:cache,mode=max \
  -t myregistry.io/myapp:latest \
  --push \
  .devcontainer/
```
