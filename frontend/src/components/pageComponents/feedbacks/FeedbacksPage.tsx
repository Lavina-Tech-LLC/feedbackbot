import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Container,
  Title,
  Paper,
  Stack,
  Text,
  Badge,
  Group,
  Card,
  SegmentedControl,
  Pagination,
  Select,
} from '@mantine/core';
import { IconLock, IconWorld } from '@tabler/icons-react';
import { useGetFeedbacks } from '@/service/feedback';
import { useGetGroups } from '@/service/group';
import type { Feedback } from '@/types';

export function FeedbacksPage() {
  const { t } = useTranslation();
  const [selectedGroup, setSelectedGroup] = useState('');
  const [filter, setFilter] = useState('all');
  const [page, setPage] = useState(1);

  const { data: groupsData } = useGetGroups('1'); // TODO: from auth context
  const groups = (groupsData as { data: { ID: number; title: string }[] } | undefined)?.data ?? [];

  const { data: feedbackData, isLoading } = useGetFeedbacks({
    groupId: selectedGroup,
    adminOnly: filter === 'admin_only' ? 'true' : filter === 'public' ? 'false' : undefined,
    page,
    limit: 20,
  });

  const feedbacks: Feedback[] =
    (feedbackData as { data: { data: Feedback[]; total: number } } | undefined)?.data?.data ?? [];
  const total =
    (feedbackData as { data: { total: number } } | undefined)?.data?.total ?? 0;
  const totalPages = Math.ceil(total / 20);

  return (
    <Container size="md" mt="xl">
      <Paper shadow="sm" p="xl" radius="md">
        <Stack>
          <Title order={2}>{t('feedbacks.title')}</Title>

          <Select
            label={t('feedbacks.selectGroup')}
            placeholder={t('feedbacks.selectGroupPlaceholder')}
            data={groups.map((g) => ({ value: String(g.ID), label: g.title }))}
            value={selectedGroup}
            onChange={(v) => {
              setSelectedGroup(v ?? '');
              setPage(1);
            }}
          />

          {selectedGroup && (
            <SegmentedControl
              value={filter}
              onChange={(v) => {
                setFilter(v);
                setPage(1);
              }}
              data={[
                { value: 'all', label: t('feedbacks.all') },
                { value: 'public', label: t('feedbacks.public') },
                { value: 'admin_only', label: t('feedbacks.adminOnly') },
              ]}
            />
          )}

          {isLoading && <Text>{t('common.loading')}</Text>}

          {!isLoading && selectedGroup && feedbacks.length === 0 && (
            <Text c="dimmed">{t('feedbacks.empty')}</Text>
          )}

          {feedbacks.map((fb) => (
            <Card key={fb.ID} shadow="xs" padding="sm" radius="sm" withBorder>
              <Group justify="space-between" mb="xs">
                <Group gap="xs">
                  {fb.admin_only ? (
                    <Badge color="red" leftSection={<IconLock size={12} />}>
                      {t('feedbacks.adminOnlyBadge')}
                    </Badge>
                  ) : (
                    <Badge color="blue" leftSection={<IconWorld size={12} />}>
                      {t('feedbacks.publicBadge')}
                    </Badge>
                  )}
                  {fb.posted && (
                    <Badge color="green" variant="light">
                      {t('feedbacks.posted')}
                    </Badge>
                  )}
                </Group>
                <Text size="xs" c="dimmed">
                  {new Date(fb.CreatedAt).toLocaleString()}
                </Text>
              </Group>
              <Text>{fb.message}</Text>
            </Card>
          ))}

          {totalPages > 1 && (
            <Pagination value={page} onChange={setPage} total={totalPages} />
          )}
        </Stack>
      </Paper>
    </Container>
  );
}
