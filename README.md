# GalaGO

Releitura de Galaga em Go usando Ebitengine.

## Instituição

- Universidade Federal do Rio Grande -- Rio Grande -- RS -- Brasil
- Engenharia da Computação

### Integrantes do projeto

- Guilherme M. Soares
- Gustavo Gimenes
- João Victor da Vara
- Lara Letittja Sague Lopez Guardiola Velloso

## Como executar

Se estiver começando do zero:

```bash
go get github.com/hajimehoshi/ebiten/v2
go mod tidy
go run .
```

Se houver erro do tipo "fatal error: X11/extensions/Xrandr.h: No such file or directory" e você estiver no Linux:

```bash
sudo apt-get install -y libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

## Como jogar

- `Enter`: inicia o jogo.
- `A/D` ou setas: move a nave para os lados.
- `Espaço`: atira.
- `P`: pausa ou continua a partida.
- `M`: liga ou desliga o som.
- `R`: reinicia após derrota ou vitória.
- `Q`: fecha o jogo.

## Funcionalidades

- Menu inicial com comandos na tela.
- Menu de seleção de naves
- Nave do jogador com movimento lateral e disparo contínuo.
- Nave de jogador com movimentos laterais, tiros duplos e mais vida.
- Meteoros aleatórios com tamanhos e resistências diferentes.
- Sistema de vidas, invulnerabilidade temporária e pontuação com multiplicador.
- Penalidade de pontos quando meteoros escapam.
- Chefe Berri com aparições periódicas, barra de vida e dificuldade progressiva.
- Ataques carregados, tiros em padrões variados e explosões.
- Power-ups temporários de tiro máximo e escudo.
- Meteoros grandes podem quebrar em fragmentos menores.
- Tela de derrota, tela de vitória, pontuação final, rank e recorde local.
- Efeitos sonoros, música ambiente e opção de mutar o som.
- Sprites carregados por `go:embed` a partir da pasta `assets/`.

## Estrutura

- `main.go`: ponto de entrada que configura a janela e inicia o loop da Ebitengine.
- `assets`: sprites reais do projeto e pacote Go com `go:embed`.
- `internal/assets`: organiza os sprites embutidos para a camada de jogo.
- `internal/config`: constantes globais, como tamanho da tela e velocidades.
- `internal/entities`: entidades do jogo, como jogador, tiros, meteoros e power-ups.
- `internal/game`: estado principal do jogo e implementação de `Update`, `Draw` e `Layout`.
