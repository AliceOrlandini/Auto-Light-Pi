import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';

class User {
  final String id;
  final String username;
  final String email;

  User({required this.id, required this.username, required this.email});

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      username: json['username'] as String,
      email: json['email'] as String,
    );
  }

  UserEntity toEntity() {
    // TODO: change this
    return const UserEntity(name: 'Mario', surname: 'Rossi');
  }
}
