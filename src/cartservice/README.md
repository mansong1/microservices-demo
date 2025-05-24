# Cart Service

The Cart service provides shopping cart functionality for the application. It manages user shopping carts and cart items.

## Building locally

The Cart service is built using .NET. To build the Cart Service locally, run:

```
dotnet build cartservice.csproj
```

To run the service:

```
dotnet run
```

## Building docker image

From `src/cartservice/`, run:

```
docker build ./
```

## Testing

To run tests:

```
dotnet test tests/
```

Or use the provided script:

```
./run-tests.sh
