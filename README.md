## Variáveis de ambiente 
|Variável|Descrição|Obrigatório|Default|
|--|--|--|--|
|REPOEAR_HOST_PORT|corresponde  aporta de exposição interna da repoear|não|8000|
|REPOEAR_DIR| correponde ao path de onde serão entregues os seus repositórios, este é mapeado sempre com um disco persistente |não|/repoear_dir/|
|GIT_SSH_PRIVATE_KEY|rresponde a chave privada do gitlab ou github, este objeto é fundamental para execução de seu processo. |sim|-|


> o acesso http publico ainda não esta sendo suportado, mas será disponibilizado no próximo release

## Apis livres
O repoear possui um conjunto de apis simplificadas para facilitar o entendimento da saúde da aplicação e execução de triggers, sedo elas:

Health Check o resultado deverá ser sempre com status code `200` com o valor `{"status":"ok"}`
```sh
curl -X GET localhost:8000/healthcheck
```
Sync forçado de repositório, este metódo possibilita ao client forçar a execuçnao de sincronização dos repositórios configurados eo resultado deverá ser sempre com status code `200` com o valor `{"status":"sync has been trigged"}`
```sh
curl -X POST localhost:8000/sync
```
> A opção de sync sempre irá realizar o pull de todas as ultimas alterações contidas no git.

## Arquio de configuração
Uma vez rodando no kubernetes seu pod irá usar este arquivo encodado em um base64 dentro de uma secred, que ao iniciar o contianer este é jogado para dentro de um disco persistente mapeado previamente via chart.

```sh
#time in seconds
refresh: 30
repositories:
  - name: release
    sync: false
    override: false
    url: git@github.com:ovrdoz/release.git
    script: |
      #!/bin/bash
      echo "hello i'm repoear.io"
```

## Usando helm versão 2
Para executar scripts usando helm version 2.x, basta iniciar o script com o comando abaixo

```sh
#!/bin/bash

HELM_VERSION 2
...
echo hello
```
> Este comando irá alternar a versão de helm entre 2 e 3 de acordo com o que você desejar, logo se vc quer executar em seu script um helm versão 2, basta adicionar a sintax no inicio de seu script `HELM_VERSION 2`, isso irá servir como um `set`


## Informações para desenvolvimento

para realizar o build da imagem utilize sempre o camando
```sh
docker build -t ovrdoz/repoear:v1.0.2 .
```
Para testar sua imagem execute passando o local da sua chave privada
```sh
docker run -v $(pwd)/config:/config -e GIT_SSH_PRIVATE_KEY=${GIT_SSH_PRIVATE_KEY}  -p 8000:8000  ovrdoz/repoear:v1.0.2
```
Caso queira saber o que esta executando dentro container use p docker exec abaixo
```sh
docker exec -it $(docker ps | grep ovrdoz | awk '{print $1}') bash
```