import { SignupForm } from '@/components/auth/signup-form';

export default function CadastroImobiliariaPage() {
  return (
    <SignupForm
      variant="standalone"
      redirectTo="http://localhost:3002/dashboard"
    />
  );
}
