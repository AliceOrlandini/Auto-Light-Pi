import 'package:auto_light_pi/core/errors/failure.dart';

class BadRequestFailure implements Failure {
  @override
  final String message;

  BadRequestFailure(this.message);
}

class UnauthorizedFailure implements Failure {
  @override
  final String message;

  UnauthorizedFailure(this.message);
}

class TimeoutFailure implements Failure {
  @override
  final String message;

  TimeoutFailure(this.message);
}

class ServerFailure implements Failure {
  @override
  final String message;

  ServerFailure(this.message);
}

class UnknownFailure implements Failure {
  @override
  final String message;

  UnknownFailure(this.message);
}
