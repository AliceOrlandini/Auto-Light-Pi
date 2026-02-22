import 'package:auto_light_pi/core/widgets/generic_input_field.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/validations/username.dart';
import 'package:auto_light_pi/features/login/presentation/widgets/username_input_field.dart';
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

  group('UsernameInputField', () {
    testWidgets('renders correctly with default state', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(const LoginState());

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const UsernameInputField(),
            ),
          ),
        ),
      );

      expect(find.byType(GenericInputField), findsOneWidget);
      expect(find.text('Username'), findsOneWidget);
    });

    testWidgets('shows error message when username is empty', (
      WidgetTester tester,
    ) async {
      final LoginState state = const LoginState(
        username: Username.dirty(''),
      ).copyWith(isUsernameValid: false);

      when(() => mockLoginBloc.state).thenReturn(state);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const UsernameInputField(),
            ),
          ),
        ),
      );

      expect(find.text('Please, insert the username'), findsOneWidget);
    });

    testWidgets('dispatches LoginUsernameChanged when text is entered', (
      WidgetTester tester,
    ) async {
      when(() => mockLoginBloc.state).thenReturn(const LoginState());

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const UsernameInputField(),
            ),
          ),
        ),
      );

      await tester.enterText(find.byType(GenericInputField), 'MarioRossi');

      verify(
        () => mockLoginBloc.add(const LoginUsernameChanged('MarioRossi')),
      ).called(1);
    });

    testWidgets('rebuilds only when username value or validity changes', (
      WidgetTester tester,
    ) async {
      final LoginState initialState = const LoginState();
      final LoginState sameUsernameState = const LoginState(
        username: Username.pure(),
      );
      final LoginState differentUsernameState = const LoginState(
        username: Username.dirty('NewUsername1'),
      ).copyWith(isUsernameValid: true);
      final LoginState differentValidityState = const LoginState(
        username: Username.dirty('NewUsername1'),
      ).copyWith(isUsernameValid: false);

      whenListen(
        mockLoginBloc,
        Stream<LoginState>.fromIterable(<LoginState>[
          initialState,
          sameUsernameState, // Should not rebuild
          differentUsernameState, // Should rebuild
          differentValidityState, // Should rebuild
        ]),
        initialState: initialState,
      );

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: BlocProvider<LoginBloc>.value(
              value: mockLoginBloc,
              child: const UsernameInputField(),
            ),
          ),
        ),
      );

      // Initial build
      expect(find.byType(GenericInputField), findsOneWidget);

      // Move to same username state - should not rebuild
      await tester.pump();
      expect(find.byType(GenericInputField), findsOneWidget);

      // Move to different username state - should rebuild
      await tester.pump();
      expect(find.byType(GenericInputField), findsOneWidget);

      // Move to different validity state - should rebuild
      await tester.pump();
      expect(find.byType(GenericInputField), findsOneWidget);
    });
  });
}
