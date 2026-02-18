import { useTranslation } from 'react-i18next';
import {
  TextInput,
  Button,
  Stack,
  Title,
  Paper,
  Container,
  Badge,
  Group,
  Text,
  ActionIcon,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { IconTrash, IconRobot } from '@tabler/icons-react';
import { useCreateBot } from '@/service';
import type { Bot } from '@/types';
import { useState } from 'react';

export function BotConfigPage() {
  const { t } = useTranslation();
  const createBot = useCreateBot();
  const [bots, setBots] = useState<Bot[]>([]);

  const form = useForm({
    initialValues: {
      token: '',
    },
    validate: {
      token: (v) => (v.trim().length < 10 ? 'Invalid token' : null),
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    // TODO: get tenant_id from context/auth once auth is integrated
    createBot.mutate(
      { tenant_id: 1, token: values.token },
      {
        onSuccess: (res) => {
          const bot = (res as { data: Bot }).data;
          setBots((prev) => [...prev, bot]);
          form.reset();
          notifications.show({
            title: t('common.success'),
            message: `Bot @${bot.bot_username} added`,
            color: 'green',
          });
        },
        onError: () => {
          notifications.show({
            title: t('common.error'),
            message: 'Failed to add bot. Check the token.',
            color: 'red',
          });
        },
      },
    );
  });

  return (
    <Container size="sm" mt="xl">
      <Stack gap="lg">
        <Paper shadow="sm" p="xl" radius="md">
          <form onSubmit={handleSubmit}>
            <Stack>
              <Title order={2}>{t('bot.title')}</Title>

              <TextInput
                label={t('bot.token')}
                description={t('bot.tokenHint')}
                placeholder="123456789:ABCDefGhIJKlmNoPQRsTUVwxyz"
                required
                {...form.getInputProps('token')}
              />

              <Button type="submit" loading={createBot.isPending}>
                {t('bot.addBot')}
              </Button>
            </Stack>
          </form>
        </Paper>

        {bots.length > 0 && (
          <Paper shadow="sm" p="xl" radius="md">
            <Stack>
              <Title order={3}>Bots</Title>
              {bots.map((bot) => (
                <Group key={bot.ID} justify="space-between">
                  <Group>
                    <IconRobot size={20} />
                    <div>
                      <Text fw={500}>@{bot.bot_username}</Text>
                      <Text size="sm" c="dimmed">
                        {bot.bot_name}
                      </Text>
                    </div>
                  </Group>
                  <Group>
                    <Badge color={bot.verified ? 'green' : 'red'}>
                      {bot.verified ? t('bot.verified') : t('bot.notVerified')}
                    </Badge>
                    <ActionIcon color="red" variant="subtle">
                      <IconTrash size={16} />
                    </ActionIcon>
                  </Group>
                </Group>
              ))}
            </Stack>
          </Paper>
        )}
      </Stack>
    </Container>
  );
}
