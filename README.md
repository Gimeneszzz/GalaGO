# GalaGO

Releitura de Galaga em Golang usando Ebitengine.

## Instituição

- Universidade Federal do Rio Grande -- Rio Grande -- RS -- Brasil
- Engenharia da Computação

### integrantes do projeto

- Guilherme M. Soares
- Gustavo Gimenes
- João Victor da Vara
- Lara Letittja Sague Lopez Guardiola Velloso


## Estrutura

- `main.go`: ponto de entrada. Configura a janela e inicia o loop da Ebitengine.
- `internal/config`: constantes globais, como tamanho da tela e velocidades.
- `assets`: sprites reais do projeto e um pequeno pacote Go com `go:embed`.
- `internal/assets`: organiza os sprites embutidos para a camada de jogo.
- `internal/entities`: entidades do jogo, como jogador, tiros e inimigos.
- `internal/game`: estado principal do jogo e implementação de `Update`, `Draw` e `Layout`.

## Como executar

Se estiver começando do zero:

```bash
go get github.com/hajimehoshi/ebiten/v2
go mod tidy
go run .
```

se houver erro do tipo "fatal error: X11/extensions/Xrandr.h: No such file or directory" e voce estiver no linux:

```bash
sudo apt-get install -y libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

## MVP implementado

- Tela horizontal `960x540` em proporcao 16:9.
- Nave do jogador com movimento lateral por `A/D` ou setas.
- Disparo continuo com `Espaco` e pequeno intervalo entre tiros.
- Meteoros aleatorios surgindo no topo da tela como obstaculos/inimigos.
- Meteoros pequenos, medios e grandes precisam de 1, 2 e 3 tiros.
- Colisao com meteoros, sistema de 3 vidas e invulnerabilidade temporaria.
- Pontuacao com multiplicador progressivo de `x01` ate `x10`.
- Penalidade de `-1000`, `-2000` ou `-3000` quando um meteoro escapa.
- Tela de game over e reinicio com a tecla `R`.
- Fundo simples com estrelas e planetas.
- Sprites carregados por `go:embed` a partir da pasta `assets/`.
- Placeholder automatico para meteoros caso a pasta seja removida sem querer.

## Proximos passos

- Colisao entre os tiros e os inimigos.
- Pontuacao e vidas.
- Animacao por sprite sheet.
- Formacoes de inimigos no estilo Galaga original.
- Possiveis efeitos de audio para tiros, explosoes e musica.
