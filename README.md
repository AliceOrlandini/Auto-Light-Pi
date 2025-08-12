# Auto Light Pi

Auto Light Pi is an open-source project designed to automate and optimize room lighting using a Raspberry Pi Pico, a light sensor, and a lamp. The system allows users to set a desired brightness level for their room; the lamp will automatically adjust its intensity to maintain that target, based on real-time sensor readings.

---

## Project Overview

**How it works:**  
- A **Raspberry Pi Pico** is connected both to a light sensor and to the lamp.
- The Pico continuously reads data from the light sensor placed in the room.
- The user sets a desired brightness percentage via a mobile app.
- The backend processes user preferences and sensor data, sending commands to the Pico.
- The Pico then adjusts the lamp's brightness to match the target value.

---

## System Architecture

The project is divided into three main components:

### 1. Mobile App

A cross-platform Flutter application that allows users to:
- Register and login.
- Set and update the desired room brightness.
- Monitor real-time sensor data and lamp status.
- Manage their account and preferences.

**Tech stack & features:**
- **Flutter** for Android, iOS, Web, and Desktop support.
- **BLoC** for state management.
- **Dio** for REST API communication.
- **Secure storage** for sensitive data.
- **Modular architecture** with clear separation of core, features, and presentation layers.

See [`mobile_app/README.md`](mobile_app/README.md) for details.

---

### 2. Backend

A Go-based backend that:
- Exposes RESTful APIs for authentication, user management, and device control.
- Handles business logic for brightness adjustment and sensor data processing.
- Manages user sessions and JWT-based authentication.
- Connects to a PostgreSQL database for persistent storage.
- Uses a clean, layered architecture (controllers, services, repositories, models, middleware).

**Tech stack & features:**
- **Go** for performance and simplicity.
- **Google Wire** for dependency injection.
- **PostgreSQL** for data storage.

See [`backend/README.md`](backend/README.md) for details.

### 3. Pico Firmware 
> The firmware for the Raspberry Pi Pico (`pico_firmware`) will be implemented in the future and is not yet included in this repository.
---

## Getting Started

1. **Set up the backend**  
   - Configure environment variables and database.
   - Run the Go server.

2. **Set up the mobile app**  
   - Install Flutter dependencies.
   - Run the app on your device or emulator.

3. **Connect the hardware**  
   - Flash the Raspberry Pi Pico with the appropriate firmware (to be implemented).
   - Connect the light sensor and lamp as per your circuit design.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.