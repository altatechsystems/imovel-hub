import { redirect } from 'next/navigation';

export default function HomePage() {
  // Redirect to login page
  // The login page will handle auth state and redirect to dashboard if already logged in
  redirect('/login');
}
