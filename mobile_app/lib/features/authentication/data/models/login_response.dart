import 'package:auto_light_pi/features/authentication/data/models/user.dart';

class LoginResponse {
  final String message;
  final User user;

  LoginResponse({required this.message, required this.user});

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      message: json['message'] as String,
      user: User.fromJson(json['user'] as Map<String, dynamic>),
    );
  }
}
