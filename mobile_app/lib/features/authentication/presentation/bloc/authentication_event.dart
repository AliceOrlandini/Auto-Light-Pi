import 'package:equatable/equatable.dart';

abstract class AuthenticationEvent extends Equatable {
  const AuthenticationEvent();
  // coverage:ignore-start
  @override
  List<Object> get props => <Object>[];
  // coverage:ignore-end
}

final class AuthCheckRequest extends AuthenticationEvent {}

final class LoginRequest extends AuthenticationEvent {
  final String username;
  final String password;

  const LoginRequest(this.username, this.password);

  @override
  List<Object> get props => <Object>[username, password];
}

final class LogoutRequest extends AuthenticationEvent {}
