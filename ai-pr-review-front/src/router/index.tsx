import { createBrowserRouter, Navigate } from 'react-router-dom';
import { MainLayout } from '@/layouts/MainLayout';
import { ReviewCreatePage } from '@/pages/ReviewCreate';
import { ReviewListPage } from '@/pages/ReviewList';
import { ReviewDetailPage } from '@/pages/ReviewDetail';
import { NotFoundPage } from '@/pages/NotFound';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: <Navigate to="/review/new" replace /> },
      { path: 'review/new', element: <ReviewCreatePage /> },
      { path: 'reviews', element: <ReviewListPage /> },
      { path: 'reviews/:id', element: <ReviewDetailPage /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
]);
