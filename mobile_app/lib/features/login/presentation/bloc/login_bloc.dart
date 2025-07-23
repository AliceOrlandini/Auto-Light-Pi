import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:auto_light_pi/features/login/presentation/validations/username.dart';
import 'package:auto_light_pi/features/login/presentation/validations/password.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:formz/formz.dart';

class LoginBloc extends Bloc<LoginEvent, LoginState> {
  final AuthenticationBloc _authenticationBloc;

  LoginBloc({required AuthenticationBloc authenticationBloc})
    : _authenticationBloc = authenticationBloc,
      super(const LoginState()) {
    on<LoginUsernameChanged>(_onUsernameChanged);
    on<LoginPasswordChanged>(_onPasswordChanged);
    on<LoginSubmitted>(_onSubmitted);
  }

  void _onUsernameChanged(
    LoginUsernameChanged event,
    Emitter<LoginState> emit,
  ) {
    final Username username = Username.dirty(event.username);
    emit(
      state.copyWith(
        username: username,
        isUsernameValid: Formz.validate(<Username>[username]),
      ),
    );
  }

  void _onPasswordChanged(
    LoginPasswordChanged event,
    Emitter<LoginState> emit,
  ) {
    final Password password = Password.dirty(event.password);
    emit(
      state.copyWith(
        password: password,
        isPasswordValid: Formz.validate(<Password>[password]),
      ),
    );
  }

  void _onSubmitted(LoginSubmitted event, Emitter<LoginState> emit) async {
    if (!state.isUsernameValid || !state.isPasswordValid) {
      return;
    }

    emit(state.copyWith(status: FormzSubmissionStatus.inProgress));

    _authenticationBloc.add(
      LoginRequest(state.username.value, state.password.value),
    );

    await emit.forEach<AuthenticationState>(
      _authenticationBloc.stream,
      onData: (AuthenticationState authenticationState) {
        if (authenticationState.status ==
            AuthenticationStatus.unauthenticated) {
          return state.copyWith(
            status: FormzSubmissionStatus.failure,
            errorMessage: authenticationState.errorMessage,
          );
        }

        if (authenticationState.status == AuthenticationStatus.authenticated) {
          return state.copyWith(status: FormzSubmissionStatus.success);
        }

        return state.copyWith(status: FormzSubmissionStatus.inProgress);
      },
    );
  }
}
