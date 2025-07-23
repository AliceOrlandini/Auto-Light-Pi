import 'package:auto_light_pi/core/errors/failure.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/login_use_case.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:dartz/dartz.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class AuthenticationBloc
    extends Bloc<AuthenticationEvent, AuthenticationState> {
  final LoginUseCase _loginUseCase;
  AuthenticationBloc({required LoginUseCase loginUseCase})
    : _loginUseCase = loginUseCase,
      super(const AuthenticationState.unknown()) {
    // on<AuthCheckRequest>(_onAuthCheckRequest);
    on<LoginRequest>(_onLoginRequest);
    // on<LogoutRequest>(_onLogoutRequest);
  }

  Future<void> _onLoginRequest(
    LoginRequest event,
    Emitter<AuthenticationState> emit,
  ) async {
    emit(const AuthenticationState.unknown());
    final Either<UserEntity, Failure> result = await _loginUseCase.call(
      username: event.username,
      password: event.password,
    );

    result.fold(
      (UserEntity user) {
        emit(AuthenticationState.authenticated(user));
      },
      (Failure failure) {
        emit(AuthenticationState.unauthenticated(failure.message));
      },
    );
  }
}
