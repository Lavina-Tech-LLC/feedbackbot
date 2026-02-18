import { useTranslation } from 'react-i18next';
import { TextInput, Button, Stack, Title, Paper, Container } from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { useCreateTenant } from '@/service';
import { useNavigate } from '@tanstack/react-router';

export function SetupPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const createTenant = useCreateTenant();

  const form = useForm({
    initialValues: {
      name: '',
      slug: '',
    },
    validate: {
      name: (v) => (v.trim().length < 2 ? 'Name is too short' : null),
      slug: (v) =>
        /^[a-z0-9-]+$/.test(v) ? null : 'Slug must be lowercase alphanumeric with dashes',
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    createTenant.mutate(values, {
      onSuccess: () => {
        notifications.show({
          title: t('common.success'),
          message: t('setup.title'),
          color: 'green',
        });
        navigate({ to: '/settings/bot' });
      },
      onError: () => {
        notifications.show({
          title: t('common.error'),
          message: 'Failed to create organization',
          color: 'red',
        });
      },
    });
  });

  return (
    <Container size="sm" mt="xl">
      <Paper shadow="sm" p="xl" radius="md">
        <form onSubmit={handleSubmit}>
          <Stack>
            <Title order={2}>{t('setup.title')}</Title>

            <TextInput
              label={t('setup.orgName')}
              placeholder="My Company"
              required
              {...form.getInputProps('name')}
            />

            <TextInput
              label={t('setup.orgSlug')}
              description={t('setup.orgSlugHint')}
              placeholder="my-company"
              required
              {...form.getInputProps('slug')}
            />

            <Button type="submit" loading={createTenant.isPending}>
              {t('setup.createOrg')}
            </Button>
          </Stack>
        </form>
      </Paper>
    </Container>
  );
}
