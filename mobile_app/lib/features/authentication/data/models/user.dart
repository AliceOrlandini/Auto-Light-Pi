import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';

class User {
  final String username;
  final String email;
  final String name;
  final String surname;

  User({
    required this.username,
    required this.email,
    required this.name,
    required this.surname,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      username: json['username'] as String,
      email: json['email'] as String,
      name: json['name'] as String,
      surname: json['surname'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      'username': username,
      'email': email,
      'name': name,
      'surname': surname,
    };
  }

  UserEntity toEntity() {
    return UserEntity(
      username: username,
      email: email,
      name: name,
      surname: surname,
    );
  }
}
