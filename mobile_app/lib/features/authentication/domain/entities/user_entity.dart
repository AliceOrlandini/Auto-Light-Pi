class UserEntity {
  final String username;
  final String email;
  final String name;
  final String surname;

  const UserEntity({
    required this.username,
    required this.email,
    required this.name,
    required this.surname,
  });

  static const UserEntity empty = UserEntity(
    username: '-',
    email: '-',
    name: '-',
    surname: '-',
  );
}
