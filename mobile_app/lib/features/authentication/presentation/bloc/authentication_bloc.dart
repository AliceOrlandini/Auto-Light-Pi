import 'package:auto_light_pi/core/failure/network_failure.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/check_authentication.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/login.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/logout.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:dartz/dartz.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class AuthenticationBloc
    extends Bloc<AuthenticationEvent, AuthenticationState> {
  final LoginUseCase _loginUseCase;
  final CheckAuthenticationUseCase _checkAuthenticationUseCase;
  final LogoutUseCase _logoutUseCase;
  AuthenticationBloc({
    required LoginUseCase loginUseCase,
    required CheckAuthenticationUseCase checkAuthenticationUseCase,
    required LogoutUseCase logoutUseCase,
  }) : _loginUseCase = loginUseCase,
       _checkAuthenticationUseCase = checkAuthenticationUseCase,
       _logoutUseCase = logoutUseCase,
       super(const AuthenticationState.unknown()) {
    on<AuthCheckRequest>(_onAuthCheckRequest);
    on<LoginRequest>(_onLoginRequest);
    on<LogoutRequest>(_onLogoutRequest);
  }

  Future<void> _onLoginRequest(
    LoginRequest event,
    Emitter<AuthenticationState> emit,
  ) async {
    emit(const AuthenticationState.unknown());
    final Either<UserEntity, NetworkFailure> result = await _loginUseCase.call(
      username: event.username,
      password: event.password,
    );

    result.fold(
      (UserEntity user) {
        emit(AuthenticationState.authenticated(user));
      },
      (NetworkFailure failure) {
        emit(
          AuthenticationState.unauthenticated(
            failure.message,
            statusCode: failure.statusCode,
          ),
        );
      },
    );
  }

  Future<void> _onAuthCheckRequest(
    AuthCheckRequest event,
    Emitter<AuthenticationState> emit,
  ) async {
    emit(const AuthenticationState.unknown());
    final UserEntity? result = await _checkAuthenticationUseCase.call();
    if (result != null) {
      emit(AuthenticationState.authenticated(result));
    } else {
      emit(const AuthenticationState.unauthenticated('Unauthenticated'));
    }
  }

  Future<void> _onLogoutRequest(
    LogoutRequest event,
    Emitter<AuthenticationState> emit,
  ) async {
    await _logoutUseCase.call();
    emit(const AuthenticationState.unauthenticated('Logged out'));
  }
}
