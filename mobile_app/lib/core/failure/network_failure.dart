import 'package:auto_light_pi/core/failure/failure.dart';

abstract class NetworkFailure extends Failure {
  final int statusCode;

  NetworkFailure(this.statusCode);

  @override
  String toString() => 'NetworkFailure($statusCode): $message';
}

class BadRequestFailure extends NetworkFailure {
  @override
  final String message;
  BadRequestFailure(this.message) : super(400);
}

class UnauthorizedFailure extends NetworkFailure {
  @override
  final String message;
  UnauthorizedFailure(this.message) : super(401);
}

class TimeoutFailure extends NetworkFailure {
  @override
  final String message;

  TimeoutFailure(this.message) : super(408);
}

class InternalServerFailure extends NetworkFailure {
  @override
  final String message;

  InternalServerFailure(this.message) : super(500);
}

class UnknownFailure extends NetworkFailure {
  @override
  final String message;

  UnknownFailure(this.message) : super(-1);
}
