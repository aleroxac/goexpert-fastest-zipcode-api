# goexpert-fatest-zipcode-api
Segundo desafio do treinamento GoExpert(FullCycle).



## O desafio
Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
As duas requisições serão feitas simultaneamente para as seguintes APIs:
- https://brasilapi.com.br/api/cep/v1/01153000 + cep
- http://viacep.com.br/ws/" + cep + "/json/



## Como rodar o projeto
``` shell
make run
```



## Funcionalidades da linguagem utilizadas
- [x] context
- [x] net/http
- [x] encoding/json
- [x] go-routines
- [x] channels
- [x] select



## Requisitos
- [x] Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- [x] O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
- [x] Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
