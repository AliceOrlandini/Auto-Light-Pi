import 'package:auto_light_pi/features/authentication/data/models/user.dart';

class LoginValidResponse {
  final String message;
  final User user;

  LoginValidResponse({required this.message, required this.user});

  factory LoginValidResponse.fromJson(Map<String, dynamic> json) {
    return LoginValidResponse(
      message: json['message'] as String,
      user: User.fromJson(json['user'] as Map<String, dynamic>),
    );
  }
}
