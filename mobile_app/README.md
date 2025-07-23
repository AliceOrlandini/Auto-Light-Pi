## Mobile App Project Structure

All Dart source files live under `lib/`, organized into the following directories:

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
Each feature follows a **clean architecture** structure:
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
