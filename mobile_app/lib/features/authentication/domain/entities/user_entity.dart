class UserEntity {
  final String name;
  final String surname;

  const UserEntity({required this.name, required this.surname});

  static const UserEntity empty = UserEntity(name: '-', surname: '-');
}
