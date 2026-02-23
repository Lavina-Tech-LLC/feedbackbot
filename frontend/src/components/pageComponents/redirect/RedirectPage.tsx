import { useEffect } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { Center, Loader } from '@mantine/core';
import { api } from '@/api';

export function RedirectPage() {
  const navigate = useNavigate();

  useEffect(() => {
    api
      .get('/auth/me')
      .then((res: any) => {
        const tenantId = res?.data?.tenant_id ?? 0;
        if (tenantId > 0) {
          navigate({ to: '/settings/bot' });
        } else {
          navigate({ to: '/setup' });
        }
      })
      .catch(() => {
        navigate({ to: '/setup' });
      });
  }, [navigate]);

  return (
    <Center h="80vh">
      <Loader />
    </Center>
  );
}
