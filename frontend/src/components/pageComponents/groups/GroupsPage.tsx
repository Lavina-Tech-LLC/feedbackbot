import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Container,
  Title,
  Paper,
  Table,
  Switch,
  Badge,
  ActionIcon,
  Modal,
  Stack,
  Checkbox,
  NumberInput,
  Button,
  Text,
  Group as MGroup,
} from '@mantine/core';
import { IconSettings } from '@tabler/icons-react';
import { useGetGroups, useUpdateGroup, useUpdateGroupConfig } from '@/service/group';
import { useCurrentTenant } from '@/utils/useCurrentTenant';
import type { Group } from '@/types';

export function GroupsPage() {
  const { t } = useTranslation();
  const tenantId = useCurrentTenant();
  const { data, isLoading } = useGetGroups(tenantId);
  const updateGroup = useUpdateGroup();
  const updateConfig = useUpdateGroupConfig();

  const [configModal, setConfigModal] = useState<Group | null>(null);
  const [postToGroup, setPostToGroup] = useState(false);
  const [forumTopicId, setForumTopicId] = useState<number | string>('');

  const groups: Group[] = (data as { data: Group[] } | undefined)?.data ?? [];

  const handleToggleActive = (group: Group) => {
    updateGroup.mutate({ id: group.ID, data: { is_active: !group.is_active } });
  };

  const openConfig = (group: Group) => {
    setConfigModal(group);
    setPostToGroup(false);
    setForumTopicId('');
  };

  const saveConfig = () => {
    if (!configModal) return;
    updateConfig.mutate(
      {
        groupId: configModal.ID,
        data: {
          post_to_group: postToGroup,
          forum_topic_id: forumTopicId ? Number(forumTopicId) : undefined,
        },
      },
      {
        onSuccess: () => setConfigModal(null),
      },
    );
  };

  return (
    <Container size="md" mt="xl">
      <Paper shadow="sm" p="xl" radius="md">
        <Title order={2} mb="md">
          {t('groups.title')}
        </Title>

        {isLoading ? (
          <Text>{t('common.loading')}</Text>
        ) : groups.length === 0 ? (
          <Text c="dimmed">{t('groups.empty')}</Text>
        ) : (
          <Table>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>{t('groups.name')}</Table.Th>
                <Table.Th>{t('groups.type')}</Table.Th>
                <Table.Th>{t('groups.status')}</Table.Th>
                <Table.Th>{t('groups.actions')}</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {groups.map((group) => (
                <Table.Tr key={group.ID}>
                  <Table.Td>{group.title}</Table.Td>
                  <Table.Td>
                    <Badge variant="light">{group.type}</Badge>
                  </Table.Td>
                  <Table.Td>
                    <Switch
                      checked={group.is_active}
                      onChange={() => handleToggleActive(group)}
                    />
                  </Table.Td>
                  <Table.Td>
                    <ActionIcon variant="subtle" onClick={() => openConfig(group)}>
                      <IconSettings size={18} />
                    </ActionIcon>
                  </Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        )}
      </Paper>

      <Modal
        opened={!!configModal}
        onClose={() => setConfigModal(null)}
        title={t('groups.configTitle')}
      >
        <Stack>
          <Text fw={500}>{configModal?.title}</Text>
          <Checkbox
            label={t('groups.postToGroup')}
            checked={postToGroup}
            onChange={(e) => setPostToGroup(e.currentTarget.checked)}
          />
          <NumberInput
            label={t('groups.forumTopicId')}
            description={t('groups.forumTopicIdHint')}
            value={forumTopicId}
            onChange={setForumTopicId}
          />
          <MGroup justify="flex-end">
            <Button variant="subtle" onClick={() => setConfigModal(null)}>
              {t('common.cancel')}
            </Button>
            <Button onClick={saveConfig} loading={updateConfig.isPending}>
              {t('common.save')}
            </Button>
          </MGroup>
        </Stack>
      </Modal>
    </Container>
  );
}
