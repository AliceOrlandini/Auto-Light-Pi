# Auto Light Pi

Auto Light Pi is a cross-platform Flutter application designed to control and monitor your lighting system.

## Features
The features implemented *until now* include:
- **User Authentication**: Secure login and registration with form validation.
- **State Management**: Uses Bloc for predictable state management.
- **Secure Storage**: Stores sensitive data securely on device.
- **REST API Integration**: Communicates with a backend server for authentication and device data.

## Getting Started

### Prerequisites

- [Flutter SDK](https://flutter.dev/docs/get-started/install) (version 3.0 or higher recommended)
- [Dart SDK](https://dart.dev/get-dart)
- A running instance of the Auto Light Pi backend (see backend documentation)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/AliceOrlandini/Auto-Light-Pi.git
   cd auto-light-pi/mobile_app
   ```

2. **Install dependencies:**
   ```bash
   flutter pub get
   ```

4. **Run the app:**
     ```bash
     flutter run
     ```

## Project Structure
All Dart source files live under `lib/`, organized into the following directories:
```
lib/
  app.dart                # App root and providers
  main.dart               # Entry point
  core/                   # Core utilities, theme, widgets, network, storage
  features/               # Feature modules (authentication, device, etc.)
  navigation/             # App routing
  di/                     # Dependency injection setup
assets/                   # Images, fonts, icons
test/                     # Unit and widget tests
```
In particular:

### 1. `core/`
Shared resources and utilities used across the entire app:
- **`core/theme/`**
  Defines color palettes, text styles, and dark/light mode support.  
- **`core/utils/`**
  Helper functions (e.g. `isValidIp()`, formatting utilities).  
- **`core/constants/`**
  Global constants (e.g. `BASE_URL`, API keys).
- **`core/errors/`**  
  Application error handling: `Failure`, `AppException`, etc.  
- **`core/widgets/`**  
  Reusable UI components (buttons, loaders, input fields).

### 2. `di/`  
Dependency injection setup using `get_it`.  
- Registers all services, repositories, data sources, and BLoC instances.

### 3. `network/`
External communication layer and HTTP interceptors:  
- API client configuration.
- Logging, authentication headers, error handling.

### 4. `navigation/`
App navigation routes and router configuration.

### 5. `features/`  
"A feature is what the user does and not what the user sees.". The components of a feature are organized into three layers:
- **`data/`**  
  Handles all data operations:  
  - **`models/`**: Define classes for JSON serialization and deserialization.  
  - **`data_sources/`**: Provide methods to fetch or persist data (using the network class).  
  - **`repositories/`**: Implement domain repository interfaces, mapping external data formats into domain models.

- **`domain/`**  
  Contains pure business rules, with no external dependencies:
  - **`entities/`**: Core domain objects (e.g. `User`, `Product`).  
  - **`repositories/`**: Interfaces that the data layer must implement.  
  - **`use_cases/`**: Orchestrate business operations (e.g. `FetchUserProfile`, `SubmitOrder`).

- **`presentation/`**  
  Manages state and UI for the feature:
  - **`bloc/`**: Implements the BLoC pattern.  
  - **`screens/`**: High-level widgets composing each screen.  
  - **`widgets/`**: Small, feature-specific UI components (custom buttons, cards, forms).

## Development

- **State Management:** [flutter_bloc](https://pub.dev/packages/flutter_bloc)
- **Dependency Injection:** [get_it](https://pub.dev/packages/get_it)
- **Networking:** [dio](https://pub.dev/packages/dio)
- **Secure Storage:** [flutter_secure_storage](https://pub.dev/packages/flutter_secure_storage)
- **Routing:** [go_router](https://pub.dev/packages/go_router)

## Running Tests

```bash
flutter test
```

## License

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.
