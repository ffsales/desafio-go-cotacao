# desafio-go-cotacao

Aplicação tem o objetivo de demonstrar o uso de uma cominicação entre aplicações client/server utilizando a linguagem GO
A aplicação server faz recebe requisições em localhost:8080/cotacao, faz uma requisição em uma api aberta de consulta de cotação de dólar, grava no banco de dados e retonar o valor atual do dólar
A aplicação client faz a requisição na aplicação server e grava o resultado em um arquivo txt
Ambas as aplicações executam as chamadas utilizando context para controlar o timeout
