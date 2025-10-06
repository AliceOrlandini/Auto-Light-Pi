import 'package:auto_light_pi/core/widgets/generic_input_field.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/validations/password.dart';
import 'package:auto_light_pi/features/login/presentation/widgets/password_input_field.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

import '../../../../mocks/mocks.dart';

void main() {
  late MockLoginBloc mockLoginBloc;

  setUp(() {
    mockLoginBloc = MockLoginBloc();
  });

  tearDown(() {
    mockLoginBloc.close();
  });

  group('PasswordInputField', () {
    testWidgets('renders correctly with default state', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(const LoginState());

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );

      expect(find.byType(GenericInputField), findsOneWidget);
      expect(find.text('Password'), findsOneWidget);
    });

    testWidgets('shows error message when password is empty', (
      WidgetTester tester,
    ) async {
      final LoginState state = const LoginState(
        password: Password.dirty(''),
      ).copyWith(isPasswordValid: false);

      when(() => mockLoginBloc.state).thenReturn(state);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );
      expect(find.text('Password cannot be empty'), findsOneWidget);
    });

    testWidgets('shows error message when password too short', (
      WidgetTester tester,
    ) async {
      final LoginState state = const LoginState(
        password: Password.dirty('a'),
      ).copyWith(isPasswordValid: false);

      when(() => mockLoginBloc.state).thenReturn(state);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );
      expect(find.text('Password too short'), findsOneWidget);
    });

    testWidgets('shows error message when password is missing a digit', (
      WidgetTester tester,
    ) async {
      final LoginState state = const LoginState(
        password: Password.dirty('Testtest'),
      ).copyWith(isPasswordValid: false);

      when(() => mockLoginBloc.state).thenReturn(state);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );
      expect(find.text('Password must contain a digit'), findsOneWidget);
    });

    testWidgets(
      'shows error message when password is missing an upper case letter',
      (WidgetTester tester) async {
        final LoginState state = const LoginState(
          password: Password.dirty('testtest123'),
        ).copyWith(isPasswordValid: false);

        when(() => mockLoginBloc.state).thenReturn(state);

        await tester.pumpWidget(
          MaterialApp(
            home: Scaffold(
              body: BlocProvider<LoginBloc>.value(
                value: mockLoginBloc,
                child: const PasswordInputField(),
              ),
            ),
          ),
        );
        expect(
          find.text('Password must contain an uppercase letter'),
          findsOneWidget,
        );
      },
    );

    testWidgets(
      'shows error message when password is missing a lower case letter',
      (WidgetTester tester) async {
        final LoginState state = const LoginState(
          password: Password.dirty('TESTTEST123'),
        ).copyWith(isPasswordValid: false);

        when(() => mockLoginBloc.state).thenReturn(state);

        await tester.pumpWidget(
          MaterialApp(
            home: Scaffold(
              body: BlocProvider<LoginBloc>.value(
                value: mockLoginBloc,
                child: const PasswordInputField(),
              ),
            ),
          ),
        );
        expect(
          find.text('Password must contain a lowercase letter'),
          findsOneWidget,
        );
      },
    );

    testWidgets('toggles obscureText when suffix icon is pressed', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(const LoginState());

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );

      // Initially, should show "visibility" icon
      expect(find.byIcon(Icons.visibility), findsOneWidget);

      // On icon tap, should toggle to "visibility_off"
      await tester.tap(find.byIcon(Icons.visibility));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.visibility_off), findsOneWidget);
    });

    testWidgets('dispatches LoginPasswordChanged when text is entered', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(const LoginState());

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );

      await tester.enterText(find.byType(GenericInputField), 'MyPassword1');

      verify(
        () => mockLoginBloc.add(const LoginPasswordChanged('MyPassword1')),
      ).called(1);
    });

    testWidgets('rebuilds only when password value or validity changes', (
      WidgetTester tester,
    ) async {
      final LoginState initialState = const LoginState();
      final LoginState samePasswordState = const LoginState(
        password: Password.pure(),
      );
      final LoginState differentPasswordState = const LoginState(
        password: Password.dirty('NewPass1'),
      ).copyWith(isPasswordValid: true);
      final LoginState differentValidityState = const LoginState(
        password: Password.dirty('NewPass1'),
      ).copyWith(isPasswordValid: false);

      whenListen(
        mockLoginBloc,
        Stream<LoginState>.fromIterable([
          initialState,
          samePasswordState, // Should not rebuild
          differentPasswordState, // Should rebuild
          differentValidityState, // Should rebuild
        ]),
        initialState: initialState,
      );

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const PasswordInputField(),
            ),
          ),
        ),
      );

      // Initial build
      expect(find.byType(GenericInputField), findsOneWidget);

      // Move to same password state - should not rebuild
      await tester.pump();
      expect(find.byType(GenericInputField), findsOneWidget);

      // Move to different password state - should rebuild
      await tester.pump();
      expect(find.byType(GenericInputField), findsOneWidget);

      // Move to different validity state - should rebuild
      await tester.pump();
      expect(find.byType(GenericInputField), findsOneWidget);
    });
  });
}
