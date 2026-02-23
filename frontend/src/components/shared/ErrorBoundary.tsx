import { Component, type ReactNode } from 'react';
import { Container, Title, Text, Button, Stack, Paper } from '@mantine/core';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return (
        <Container size="sm" mt="xl">
          <Paper shadow="sm" p="xl" radius="md">
            <Stack align="center" gap="md">
              <Title order={2}>Something went wrong</Title>
              <Text c="dimmed">{this.state.error?.message}</Text>
              <Button onClick={() => window.location.reload()}>Reload page</Button>
            </Stack>
          </Paper>
        </Container>
      );
    }

    return this.props.children;
  }
}
