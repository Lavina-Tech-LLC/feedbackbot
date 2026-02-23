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
  Loader,
  Center,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { IconTrash, IconRobot } from '@tabler/icons-react';
import { useGetBots, useCreateBot, useDeleteBot } from '@/service';
import { useCurrentTenant } from '@/utils/useCurrentTenant';
import type { Bot } from '@/types';

export function BotConfigPage() {
  const { t } = useTranslation();
  const tenantId = useCurrentTenant();
  const createBot = useCreateBot();
  const deleteBot = useDeleteBot();
  const { data: botsData, isLoading } = useGetBots();
  const bots: Bot[] = (botsData as { data: Bot[] } | undefined)?.data ?? [];

  const form = useForm({
    initialValues: {
      token: '',
    },
    validate: {
      token: (v) => (v.trim().length < 10 ? 'Invalid token' : null),
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    createBot.mutate(
      { tenant_id: Number(tenantId), token: values.token },
      {
        onSuccess: (res) => {
          const bot = (res as { data: Bot }).data;
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

  const handleDelete = (bot: Bot) => {
    deleteBot.mutate(bot.ID, {
      onSuccess: () => {
        notifications.show({
          title: t('common.success'),
          message: `Bot @${bot.bot_username} removed`,
          color: 'green',
        });
      },
    });
  };

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

        <Paper shadow="sm" p="xl" radius="md">
          <Stack>
            <Title order={3}>Bots</Title>
            {isLoading ? (
              <Center>
                <Loader size="sm" />
              </Center>
            ) : bots.length === 0 ? (
              <Text c="dimmed">No bots configured yet.</Text>
            ) : (
              bots.map((bot) => (
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
                    <ActionIcon
                      color="red"
                      variant="subtle"
                      onClick={() => handleDelete(bot)}
                      loading={deleteBot.isPending}
                    >
                      <IconTrash size={16} />
                    </ActionIcon>
                  </Group>
                </Group>
              ))
            )}
          </Stack>
        </Paper>
      </Stack>
    </Container>
  );
}
