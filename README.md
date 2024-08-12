# Prova de Suficiência WEB II - 2024/2

O projeto realizado é uma API de gerenciamento de comandas, utilizando o padrão MVC. A API permite operações como criar, atualizar, deletar e listar comandas. A autenticação é realizada via tokens JWT (JSON Web Token) nas rotas POST, PUT e DELETE.

## Tecnologias utilizadas :computer:

- **Linguagem de Programação**: Go (Golang)
- **Framework Web**: [Gin](https://github.com/gin-gonic/gin)
- **Banco de Dados**: MongoDB
- **ORM**: [MongoDB Driver for Go](https://github.com/mongodb/mongo-go-driver)
- **Autenticação**: [JWT GO](https://github.com/dgrijalva/jwt-go)

## Requisitos :books:

Para rodar este projeto, você precisa ter instalado:

- Go 1.20 ou superior
- MongoDB (servidor)
- Git (para clonar o repositório)
- Dependências do projeto (gerenciadas com `go mod`)

## Configuração do ambiente para rodar a aplicação :mag:

### 1. Clonando o repositório

```bash
git clone https://github.com/sanarocha/prova-suficiencia-web-ii.git
cd prova-suficiencia-web-ii
```

### 2. Instalando todas as dependências

```bash
go mod tidy
```

### 3. Conectando ao BD

É necessário que o MondoDB esteja instalado e configurado na máquina na URL padrão 'mongodb://localhost:27017'.

### 4. Executando a aplicação

```bash
go run main.go
```

### 5. Acessando a documentação

A documentação da API está disponível no Swagger, depois de iniciar a aplicação, é possível acessar ela na URL abaixo.

```bash
http://localhost:8080/swagger/index.html
```

Obrigada desde já!

