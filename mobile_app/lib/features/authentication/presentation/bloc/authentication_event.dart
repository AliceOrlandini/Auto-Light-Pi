import 'package:equatable/equatable.dart';

abstract class AuthenticationEvent extends Equatable {
  const AuthenticationEvent();
  @override
  List<Object> get props => <Object>[];
}

final class AuthCheckRequest extends AuthenticationEvent {}

final class LoginRequest extends AuthenticationEvent {
  final String username;
  final String password;

  const LoginRequest(this.username, this.password);
}

final class LogoutRequest extends AuthenticationEvent {}
