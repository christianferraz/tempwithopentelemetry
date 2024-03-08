
# API de Consulta de CEP e Clima

Este projeto é uma aplicação Go simples que fornece uma API para consultar informações de CEP (Código de Endereçamento Postal) e, com base na localidade do CEP, buscar informações climáticas atuais para essa localidade utilizando a API externa WeatherAPI.

## Funcionalidades

* Consulta de informações de CEP.
* Consulta de informações climáticas baseadas na localidade do CEP.

## Como Executar

### Pré-Requisitos

* Go instalado em sua máquina (versão 1.22 ou superior recomendada).
* Acesso à internet para as consultas às APIs externas.

### Instruções

1. Clone o repositório para a sua máquina local.
2. Navegue até o diretório do projeto através do terminal.
3. Execute o comando `go run main.go` para iniciar o servidor.
4. O servidor estará disponível em `https://temp-tiqlkxh6za-uc.a.run.app`.

## Endpoints

### Consulta de CEP

* **URL** : `/cep/{cep}`
* **Método** : `GET`
* **URL Params** : Substitua `{cep}` pelo CEP desejado.
* **Resposta de Sucesso** :
* **Código** : 200 OK
* **Conteúdo** : Informações climáticas da localidade do CEP.
* **Resposta de Erro** :
* **Código** : 422 Unprocessable Entity (CEP inválido)
* **Código** : 404 Not Found (CEP não encontrado)

### Exemplo de Uso

Para consultar informações climáticas para o CEP 79052564, você acessaria:

<pre><div class="dark bg-gray-950 rounded-md"><div class="flex items-center relative text-token-text-secondary bg-token-main-surface-secondary px-4 py-2 text-xs font-sans justify-between rounded-t-md"><span>bash</span><span class="" data-state="closed"><button class="flex gap-1 items-center"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="icon-sm"><path fill-rule="evenodd" clip-rule="evenodd" d="M12 4C10.8954 4 10 4.89543 10 6H14C14 4.89543 13.1046 4 12 4ZM8.53513 4C9.22675 2.8044 10.5194 2 12 2C13.4806 2 14.7733 2.8044 15.4649 4H17C18.6569 4 20 5.34315 20 7V19C20 20.6569 18.6569 22 17 22H7C5.34315 22 4 20.6569 4 19V7C4 5.34315 5.34315 4 7 4H8.53513ZM8 6H7C6.44772 6 6 6.44772 6 7V19C6 19.5523 6.44772 20 7 20H17C17.5523 20 18 19.5523 18 19V7C18 6.44772 17.5523 6 17 6H16C16 7.10457 15.1046 8 14 8H10C8.89543 8 8 7.10457 8 6Z" fill="currentColor"></path></svg>Copy code</button></span></div><div class="p-4 overflow-y-auto"><code class="!whitespace-pre hljs language-bash">GET https://temp-tiqlkxh6za-uc.a.run.app/cep/79052564
</code></div></div></pre>

## Notas Adicionais

* Este projeto utiliza a API Viacep ([https://viacep.com.br](https://viacep.com.br/)) para consulta de informações de CEP e a API WeatherAPI ([http://api.weatherapi.com](http://api.weatherapi.com/)) para informações climáticas.
