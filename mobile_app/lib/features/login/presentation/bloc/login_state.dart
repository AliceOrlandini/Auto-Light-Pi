import 'package:equatable/equatable.dart';
import 'package:formz/formz.dart';
import 'package:auto_light_pi/features/login/presentation/validations/username.dart';
import 'package:auto_light_pi/features/login/presentation/validations/password.dart';

final class LoginState extends Equatable {
  final Username username;
  final Password password;
  final bool isUsernameValid;
  final bool isPasswordValid;
  final FormzSubmissionStatus status;
  final String? errorMessage;

  const LoginState({
    this.username = const Username.pure(),
    this.password = const Password.pure(),
    this.isUsernameValid = false,
    this.isPasswordValid = false,
    this.status = FormzSubmissionStatus.initial,
    this.errorMessage,
  });

  LoginState copyWith({
    FormzSubmissionStatus? status,
    Username? username,
    Password? password,
    bool? isUsernameValid,
    bool? isPasswordValid,
    String? errorMessage,
  }) {
    return LoginState(
      status: status ?? this.status,
      username: username ?? this.username,
      password: password ?? this.password,
      isUsernameValid: isUsernameValid ?? this.isUsernameValid,
      isPasswordValid: isPasswordValid ?? this.isPasswordValid,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  @override
  List<Object?> get props => <Object?>[
    status,
    username,
    password,
    isUsernameValid,
    isPasswordValid,
    errorMessage,
  ];
}
