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
  TextInput,
  Button,
} from '@mantine/core';
import { useDebouncedValue } from '@mantine/hooks';
import { DatePickerInput } from '@mantine/dates';
import { IconLock, IconWorld, IconCalendar, IconSearch, IconDownload } from '@tabler/icons-react';
import { useGetFeedbacks, getExportCsvUrl } from '@/service/feedback';
import { useGetGroups } from '@/service/group';
import type { Feedback } from '@/types';
import '@mantine/dates/styles.css';

export function FeedbacksPage() {
  const { t } = useTranslation();
  const [selectedGroup, setSelectedGroup] = useState('');
  const [filter, setFilter] = useState('all');
  const [page, setPage] = useState(1);
  const [dateRange, setDateRange] = useState<[Date | string | null, Date | string | null]>([null, null]);
  const [search, setSearch] = useState('');
  const [debouncedSearch] = useDebouncedValue(search, 300);

  const { data: groupsData } = useGetGroups();
  const groups = (groupsData as { data: { id: number; title: string }[] } | undefined)?.data ?? [];

  const feedbackParams = {
    groupId: selectedGroup,
    adminOnly: filter === 'admin_only' ? 'true' : filter === 'public' ? 'false' : undefined,
    page,
    limit: 20,
    dateFrom: dateRange[0] ? new Date(dateRange[0]).toISOString() : undefined,
    dateTo: dateRange[1] ? new Date(dateRange[1]).toISOString() : undefined,
    search: debouncedSearch || undefined,
  };

  const { data: feedbackData, isLoading } = useGetFeedbacks(feedbackParams);

  const feedbacks: Feedback[] =
    (feedbackData as { data: { data: Feedback[]; total: number } } | undefined)?.data?.data ?? [];
  const total = (feedbackData as { data: { total: number } } | undefined)?.data?.total ?? 0;
  const totalPages = Math.ceil(total / 20);

  const handleExport = () => {
    const url = getExportCsvUrl(feedbackParams);
    window.open(url, '_blank');
  };

  return (
    <Container size="md" mt="xl">
      <Paper shadow="sm" p="xl" radius="md">
        <Stack>
          <Group justify="space-between">
            <Title order={2}>{t('feedbacks.title')}</Title>
            {selectedGroup && (
              <Button
                variant="light"
                leftSection={<IconDownload size={16} />}
                onClick={handleExport}
              >
                {t('feedbacks.exportCsv')}
              </Button>
            )}
          </Group>

          <Group grow>
            <Select
              label={t('feedbacks.selectGroup')}
              placeholder={t('feedbacks.selectGroupPlaceholder')}
              data={groups.map((g) => ({ value: String(g.id), label: g.title }))}
              value={selectedGroup}
              onChange={(v) => {
                setSelectedGroup(v ?? '');
                setPage(1);
              }}
            />

            <DatePickerInput
              type="range"
              label={t('feedbacks.dateRange')}
              placeholder={t('feedbacks.dateRangePlaceholder')}
              value={dateRange}
              onChange={(v) => {
                setDateRange(v);
                setPage(1);
              }}
              clearable
              leftSection={<IconCalendar size={16} />}
            />
          </Group>

          {selectedGroup && (
            <TextInput
              placeholder={t('feedbacks.searchPlaceholder')}
              leftSection={<IconSearch size={16} />}
              value={search}
              onChange={(e) => {
                setSearch(e.currentTarget.value);
                setPage(1);
              }}
            />
          )}

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
            <Card key={fb.id} shadow="xs" padding="sm" radius="sm" withBorder>
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
                  {new Date(fb.createdAt).toLocaleString()}
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
