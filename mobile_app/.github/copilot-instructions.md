# Auto Light Pi - AI Coding Instructions

You are working on a **Flutter** mobile application using **Clean Architecture** and **Feature-First** organization.

## Architecture & Structure
- **Pattern**: Clean Architecture partitioned by feature.
- **Layers**:
  - `presentation`: `bloc` (State Management), `screens` (UI).
  - `domain`: `use_cases`, `repositories` (Interfaces), `entities` (Business Objects).
  - `data`: `data_sources` (Remote/Local), `repositories` (Implementations), `models` (DTOs).
- **Dependency Injection**: Use `get_it` with manual registration in [lib/di/di.dart](lib/di/di.dart).
- **State Management**: `flutter_bloc`. All Blocs must extend `Bloc` or `Cubit`.
- **Navigation**: `go_router` with redirection logic based on `AuthenticationBloc` state stream ([lib/navigation/routes.dart](lib/navigation/routes.dart)).

## Core Conventions
- **Error Handling**: Use `dartz`'s `Either<L, R>` for Repository methods. `Left` is always a failure (e.g., `NetworkFailure`), `Right` is success data.
  - Example: `Future<Either<UserEntity, NetworkFailure>> authenticate(...)`
- **Networking**: `dio` is the HTTP client. Use `DioClient` wrapper.
  - Auth: `JwtTokenInterceptor` manages token injection and refresh.
- **Models**: Manual JSON serialization (`fromJson`/`toJson`). Avoid code generation libraries unless specified.
- **Equality**: Use `equatable` for States, Events, and Entities to ensure value equality.
- **Environment**: Access variables via `flutter_dotenv` (e.g., `dotenv.get('BACKEND_URL')`).

## Developing New Features
1.  **Domain First**: Define `Entity`, `Repository` interface, and `UseCase`.
2.  **Data Layer**: Implement `DataSource` (Remote/Local), `Model` (DTO), and `RepositoryImpl`.
3.  **DI Registration**: Register all new classes in `lib/di/di.dart` (Follow the existing grouping: Domain, Data, Presentation).
4.  **Presentation**: Create `Bloc`/`Cubit` and `Screen`.
5.  **Route**: Add entry to `lib/navigation/routes.dart`.

## Testing
- **Unit Tests**: Use `test` and `mocktail`.
  - Mock repositories/use-cases.
- **Bloc Tests**: Use `bloc_test`.
- **Naming**: `test/features/<feature>/...` mirroring `lib/` structure.

## Key Files
- **Entry Point**: [lib/main.dart](lib/main.dart) (Initializes `dotenv`, `DI`).
- **App Wrapper**: [lib/app.dart](lib/app.dart) (Provides global Blocs like `AuthenticationBloc`).
- **DI Setup**: [lib/di/di.dart](lib/di/di.dart).
- **Routes**: [lib/navigation/routes.dart](lib/navigation/routes.dart).
