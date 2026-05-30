import { RouterProvider } from 'react-router-dom';
import { router } from '@/router';
import { useHealth } from '@/hooks/useHealth';

export function App() {
  useHealth();
  return <RouterProvider router={router} />;
}
