# API REST de comercio electrónico en Go

## Instalación

Hay algunas herramientas que necesitas instalar para ejecutar el proyecto. Asegúrate de tener las siguientes herramientas instaladas en tu máquina.

- [Migrate (para migraciones de base de datos)](https://github.com/golang-migrate/migrate/tree/v4.17.0/cmd/migrate)

## Ejecución del proyecto

Primero asegúrate de tener una base de datos MySQL en ejecución en tu máquina o simplemente cámbiala por cualquier otro almacenamiento que prefieras bajo `/db`.

Luego crea una base de datos con el nombre que desees (*`ecom` es el predeterminado*) y ejecuta las migraciones.

```bash
make migrate-up

Luego de eso, podes correr el proyecto ejecutando el siguiente comando:

```bash
make run
```

## Running the tests

To run the tests, you can use the following command:

```bash
make test
```
