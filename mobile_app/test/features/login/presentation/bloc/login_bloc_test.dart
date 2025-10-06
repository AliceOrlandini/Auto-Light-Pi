import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/validations/password.dart';
import 'package:auto_light_pi/features/login/presentation/validations/username.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:formz/formz.dart';
import 'package:mocktail/mocktail.dart';

import '../../../../mocks/mocks.dart';

void main() {
  late MockAuthenticationBloc mockAuthenticationBloc;

  setUp(() {
    mockAuthenticationBloc = MockAuthenticationBloc();
  });

  group('LoginBloc', () {
    test('initial state of username is a pure element', () {
      expect(
        LoginBloc(authenticationBloc: mockAuthenticationBloc).state.username,
        const Username.pure(),
      );
    });
    test('initial state of isUsernameValid is false', () {
      expect(
        LoginBloc(
          authenticationBloc: mockAuthenticationBloc,
        ).state.isUsernameValid,
        false,
      );
    });

    test('initial state of form status is initial', () {
      expect(
        LoginBloc(authenticationBloc: mockAuthenticationBloc).state.status,
        FormzSubmissionStatus.initial,
      );
    });

    blocTest<LoginBloc, LoginState>(
      'emits [LoginState with an updated valid username] when LoginUsernameChanged with a valid username is added',
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      act: (LoginBloc bloc) => bloc.add(const LoginUsernameChanged('Alice')),
      expect: () => <LoginState>[
        const LoginState(
          username: Username.dirty('Alice'),
          isUsernameValid: true,
        ),
      ],
    );

    blocTest<LoginBloc, LoginState>(
      'emits [LoginState with an updated valid password] when LoginPasswordChanged with a valid password is added',
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      act: (LoginBloc bloc) =>
          bloc.add(const LoginPasswordChanged('TestPassword123')),
      expect: () => <LoginState>[
        const LoginState(
          password: Password.dirty('TestPassword123'),
          isPasswordValid: true,
        ),
      ],
    );

    blocTest<LoginBloc, LoginState>(
      'emits [LoginState with an invalid password] when LoginPasswordChanged with a password too short is added',
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      act: (LoginBloc bloc) => bloc.add(const LoginPasswordChanged('Test123')),
      expect: () => <LoginState>[
        const LoginState(
          password: Password.dirty('Test123'),
          isPasswordValid: false,
        ),
      ],
      verify: (LoginBloc bloc) {
        expect(bloc.state.password.error, PasswordValidationError.tooShort);
      },
    );

    blocTest<LoginBloc, LoginState>(
      'emits [LoginState with an invalid password] when LoginPasswordChanged with a password without a digit is added',
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      act: (LoginBloc bloc) => bloc.add(const LoginPasswordChanged('TestTest')),
      expect: () => <LoginState>[
        const LoginState(
          password: Password.dirty('TestTest'),
          isPasswordValid: false,
        ),
      ],
      verify: (LoginBloc bloc) {
        expect(bloc.state.password.error, PasswordValidationError.digitMissing);
      },
    );

    blocTest<LoginBloc, LoginState>(
      'emits [LoginState with an invalid password] when LoginPasswordChanged with a password without an upper case letter is added',
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      act: (LoginBloc bloc) =>
          bloc.add(const LoginPasswordChanged('testtest123')),
      expect: () => <LoginState>[
        const LoginState(
          password: Password.dirty('testtest123'),
          isPasswordValid: false,
        ),
      ],
      verify: (LoginBloc bloc) {
        expect(
          bloc.state.password.error,
          PasswordValidationError.upperCaseMissing,
        );
      },
    );

    blocTest<LoginBloc, LoginState>(
      'emits [LoginState with an invalid password] when LoginPasswordChanged with a password without a lower case letter is added',
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      act: (LoginBloc bloc) =>
          bloc.add(const LoginPasswordChanged('TESTTEST123')),
      expect: () => <LoginState>[
        const LoginState(
          password: Password.dirty('TESTTEST123'),
          isPasswordValid: false,
        ),
      ],
      verify: (LoginBloc bloc) {
        expect(
          bloc.state.password.error,
          PasswordValidationError.lowerCaseMissing,
        );
      },
    );

    blocTest<LoginBloc, LoginState>(
      'emits [no new state] when LoginSubmitted is called but the form is not valid',
      build: () => LoginBloc(
        authenticationBloc: mockAuthenticationBloc,
      ), // at the beginning both username and password are empty so the form is not valid
      act: (LoginBloc bloc) => bloc.add(const LoginSubmitted()),
      expect: () => <LoginState>[], // we expect nothing changes
    );

    blocTest<LoginBloc, LoginState>(
      'emits [inProgress, success] when LoginSubmitted is called and the credentials are valid',
      setUp: () {
        // pre-configure the mock to return a valid login
        whenListen(
          mockAuthenticationBloc,
          Stream<AuthenticationState>.fromIterable(<AuthenticationState>[
            const AuthenticationState.authenticated(
              UserEntity(
                username: 'alice',
                email: 'alice@gmail.com',
                name: 'Alice',
                surname: 'Smith',
              ),
            ),
          ]),
          initialState: const AuthenticationState.unknown(),
        );
      },
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      // the seed simulates the insert of valid data inside the form
      seed: () => const LoginState(
        username: Username.dirty('Alice'),
        password: Password.dirty('TestTest123'),
        isUsernameValid: true,
        isPasswordValid: true,
        status: FormzSubmissionStatus.initial,
      ),
      act: (LoginBloc bloc) => bloc.add(const LoginSubmitted()),
      expect: () => <LoginState>[
        const LoginState(
          username: Username.dirty('Alice'),
          password: Password.dirty('TestTest123'),
          isUsernameValid: true,
          isPasswordValid: true,
          status: FormzSubmissionStatus.inProgress,
        ),
        const LoginState(
          username: Username.dirty('Alice'),
          password: Password.dirty('TestTest123'),
          isUsernameValid: true,
          isPasswordValid: true,
          status: FormzSubmissionStatus.success,
        ),
      ],
      verify: (_) {
        // verify that the AuthenticationBloc LoginRequest function is actually called
        verify(
          () => mockAuthenticationBloc.add(
            const LoginRequest('Alice', 'TestTest123'),
          ),
        ).called(1);
      },
    );

    blocTest<LoginBloc, LoginState>(
      'emits [inProgress, faliure] when LoginSubmitted is called but the credentials are not valid',
      setUp: () {
        // pre-configure the mock to return a valid login
        whenListen(
          mockAuthenticationBloc,
          Stream<AuthenticationState>.fromIterable(<AuthenticationState>[
            const AuthenticationState.unauthenticated(
              'Credentials not valid',
              statusCode: 401,
            ),
          ]),
          initialState: const AuthenticationState.unknown(),
        );
      },
      build: () => LoginBloc(authenticationBloc: mockAuthenticationBloc),
      // the seed simulates the insert of valid data inside the form
      seed: () => const LoginState(
        username: Username.dirty('Alice'),
        password: Password.dirty('TestTest123'),
        isUsernameValid: true,
        isPasswordValid: true,
        status: FormzSubmissionStatus.initial,
      ),
      act: (LoginBloc bloc) => bloc.add(const LoginSubmitted()),
      expect: () => <LoginState>[
        const LoginState(
          username: Username.dirty('Alice'),
          password: Password.dirty('TestTest123'),
          isUsernameValid: true,
          isPasswordValid: true,
          status: FormzSubmissionStatus.inProgress,
        ),
        const LoginState(
          username: Username.dirty('Alice'),
          password: Password.dirty('TestTest123'),
          isUsernameValid: true,
          isPasswordValid: true,
          status: FormzSubmissionStatus.failure,
          errorMessage: 'Credentials not valid (401)',
        ),
      ],
      verify: (_) {
        // verify that the AuthenticationBloc LoginRequest function is actually called
        verify(
          () => mockAuthenticationBloc.add(
            const LoginRequest('Alice', 'TestTest123'),
          ),
        ).called(1);
      },
    );

    blocTest<LoginBloc, LoginState>(
      'emits [inProgress] when AuthenticationBloc emits a non-terminal state and _onSubmitted returns state',
      build: () {
        mockAuthenticationBloc = MockAuthenticationBloc();

        when(
          () => mockAuthenticationBloc.state,
        ).thenReturn(const AuthenticationState.unknown());

        whenListen<AuthenticationState>(
          mockAuthenticationBloc,
          Stream<AuthenticationState>.value(
            const AuthenticationState.unknown(),
          ),
        );

        return LoginBloc(authenticationBloc: mockAuthenticationBloc);
      },
      seed: () => const LoginState(
        username: Username.dirty('Alice'),
        password: Password.dirty('TestTest123'),
        isUsernameValid: true,
        isPasswordValid: true,
      ),
      act: (LoginBloc bloc) => bloc.add(const LoginSubmitted()),
      expect: () => <LoginState>[
        const LoginState(
          username: Username.dirty('Alice'),
          password: Password.dirty('TestTest123'),
          isUsernameValid: true,
          isPasswordValid: true,
          status: FormzSubmissionStatus.inProgress,
        ),
      ],
    );
  });

  group('LoginEvent', () {
    test('LoginUsernameChanged supports value equality', () {
      expect(
        const LoginUsernameChanged('Alice'),
        equals(const LoginUsernameChanged('Alice')),
      );

      expect(
        const LoginUsernameChanged('Alice'),
        isNot(const LoginUsernameChanged('Bob')),
      );
    });
  });

  group('LoginState', () {
    test(
      'should be equal when all properties match, disegual when at least one property does not match',
      () {
        expect(
          const LoginState(
            username: Username.dirty('Alice'),
            isUsernameValid: true,
          ),
          equals(
            const LoginState(
              username: Username.dirty('Alice'),
              isUsernameValid: true,
            ),
          ),
        );

        expect(
          const LoginState(
            username: Username.dirty('Alice'),
            isUsernameValid: true,
          ),
          isNot(
            const LoginState(
              username: Username.dirty('Bob'),
              isUsernameValid: true,
            ),
          ),
        );

        expect(
          const LoginState(
            password: Password.dirty('TestTest123'),
            isPasswordValid: true,
          ),
          equals(
            const LoginState(
              password: Password.dirty('TestTest123'),
              isPasswordValid: true,
            ),
          ),
        );

        expect(
          const LoginState(
            password: Password.dirty('TestTest123'),
            isPasswordValid: true,
          ),
          isNot(
            const LoginState(
              password: Password.dirty('Another123'),
              isPasswordValid: true,
            ),
          ),
        );

        expect(
          const LoginPasswordChanged('TestTest123'),
          equals(const LoginPasswordChanged('TestTest123')),
        );

        expect(
          const LoginPasswordChanged('TestTest123'),
          isNot(const LoginPasswordChanged('Another123')),
        );
      },
    );
  });
}
