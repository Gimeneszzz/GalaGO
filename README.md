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

- Tela `800x600`.
- Menu inicial com `Enter` para comecar e comandos na tela.
- Nave do jogador com movimento lateral por `A/D` ou setas.
- Disparo continuo com `Espaco` e pequeno intervalo entre tiros.
- Meteoros aleatorios surgindo no topo da tela como obstaculos/inimigos.
- Meteoros pequenos, medios e grandes precisam de 1, 2 e 3 tiros.
- Colisao com meteoros, sistema de 3 vidas e invulnerabilidade temporaria.
- Pontuacao com multiplicador progressivo de `x01` ate `x10`.
- Penalidade de `-1000`, `-2000` ou `-3000` quando um meteoro escapa.
- Tela de game over e reinicio com a tecla `R`.
- Tela de vitoria ao derrotar Berri 4 vezes.
- Pausa com a tecla `P` e recorde local da sessao.
- Encerramento do jogo com a tecla `Q`.
- Margem lateral de 20 pixels para impedir meteoros inalcançaveis.
- Fundo simples com estrelas e planetas.
- Barra de vida, aviso de aproximacao, tiro carregado e explosoes da Berri.
- A mira vermelha da Berri e um aviso; o dano vem do projetil disparado depois.
- Dificuldade progressiva conforme pontos, derrotas e aparicoes da Berri.
- Entrada animada, barra da proxima aparicao e mensagem de Berri repelida.
- Berri fica protegida durante a entrada, entao so recebe dano ao chegar na posicao.
- Power-ups temporarios de tiro duplo e escudo ao destruir meteoros.
- Feedback de dano na Berri com texto flutuante e flash quando o jogador sofre dano.
- Padroes de ataque da Berri mudam por nivel, com tiros em leque nos niveis altos.
- Alguns meteoros grandes quebram em fragmentos menores.
- Tela de vitoria com pontuacao final e rank.
- Alternancia de som com a tecla `M`.
- Efeitos sonoros simples e musica ambiente gerada pelo jogo.
- Sprites carregados por `go:embed` a partir da pasta `assets/`.
- Placeholder automatico para meteoros caso a pasta seja removida sem querer.

## Colocar o rosto do professor nos planetas

Coloque uma imagem PNG do professor em `assets/professor.png` e rode o jogo
normalmente:

```bash
go run .
```

Se esse arquivo existir, o jogo aplica automaticamente o rosto em formato
circular nos planetas do fundo. Se o arquivo nao existir, os planetas continuam
aparecendo com as sprites originais.

## Adicionar a nave do professor

Coloque a sprite da nave do professor em `assets/professornave.png`. Quando esse
arquivo existir, uma nave especial aparece alguns segundos depois do inicio da
partida e volta de tempos em tempos. Ela fica acima da nave do jogador,
acompanha o movimento horizontal, pode ser atingida pelos tiros do jogador e
dispara lasers em direcoes levemente aleatorias tentando acertar o jogador. A
cada aparicao, Berri ganha bastante vida, fica mais tempo em tela, se move mais
rapido, entra mais depressa e atira com mais pressao. Se a imagem vier com fundo
branco ou margem grande, o jogo remove esse fundo, recorta a nave e reduz o
sprite automaticamente.

## Proximos passos

- Colisao entre os tiros e os inimigos.
- Pontuacao e vidas.
- Animacao por sprite sheet.
- Formacoes de inimigos no estilo Galaga original.
- Possiveis efeitos de audio para tiros, explosoes e musica.
