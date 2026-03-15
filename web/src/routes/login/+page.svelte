<script lang="ts">
  import { goto } from '$app/navigation';
  import { api } from '$lib/api';
  import { auth } from '$lib/auth.svelte';
  import { t } from '$lib/i18n.svelte';
  import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';

  let mode = $state<'login' | 'forgot' | 'reset'>('login');
  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);
  let successMsg = $state('');

  // Forgot password state
  let resetCode = $state('');
  let newPassword = $state('');
  let confirmNewPassword = $state('');
  let cooldown = $state(0);
  let cooldownTimer: ReturnType<typeof setInterval> | null = null;

  const passwordTooShort = $derived(newPassword.length > 0 && newPassword.length < 6);
  const passwordMismatch = $derived(confirmNewPassword.length > 0 && newPassword !== confirmNewPassword);

  function startCooldown() {
    cooldown = 60;
    if (cooldownTimer) clearInterval(cooldownTimer);
    cooldownTimer = setInterval(() => {
      cooldown--;
      if (cooldown <= 0) {
        cooldown = 0;
        if (cooldownTimer) { clearInterval(cooldownTimer); cooldownTimer = null; }
      }
    }, 1000);
  }

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    loading = true;
    try {
      const res = await api.auth.login(email, password);
      auth.login(res.token);
      goto('/');
    } catch (err) {
      error = err instanceof Error ? err.message : t('login.failed');
    } finally {
      loading = false;
    }
  }

  async function handleSendResetCode(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    loading = true;
    try {
      await api.auth.forgotPassword(email);
      startCooldown();
      mode = 'reset';
    } catch (err) {
      error = err instanceof Error ? err.message : t('forgot.sendFailed');
    } finally {
      loading = false;
    }
  }

  async function handleResendCode() {
    if (cooldown > 0 || loading) return;
    error = '';
    loading = true;
    try {
      await api.auth.forgotPassword(email);
      startCooldown();
      resetCode = '';
    } catch (err) {
      error = err instanceof Error ? err.message : t('forgot.sendFailed');
    } finally {
      loading = false;
    }
  }

  async function handleResetPassword(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    if (newPassword !== confirmNewPassword) {
      error = t('forgot.passwordMismatch');
      return;
    }
    loading = true;
    try {
      await api.auth.resetPassword(email, resetCode, newPassword);
      successMsg = t('forgot.success');
      mode = 'login';
      resetCode = '';
      newPassword = '';
      confirmNewPassword = '';
    } catch (err) {
      error = err instanceof Error ? err.message : t('forgot.resetFailed');
    } finally {
      loading = false;
    }
  }

  function goToForgot() {
    error = '';
    successMsg = '';
    mode = 'forgot';
  }

  function backToLogin() {
    error = '';
    successMsg = '';
    resetCode = '';
    newPassword = '';
    confirmNewPassword = '';
    mode = 'login';
  }
</script>

<div class="min-h-screen flex items-center justify-center relative">
  <div class="absolute top-4 right-4">
    <LanguageSwitcher />
  </div>
  <div class="w-full max-w-sm space-y-6">
    <div class="text-center">
      <h1 class="text-2xl font-bold">{t('login.title')}</h1>
      <p class="text-sm text-muted-foreground mt-1">
        {#if mode === 'login'}
          {t('login.subtitle')}
        {:else}
          {mode === 'forgot' ? t('forgot.subtitle') : t('forgot.title')}
        {/if}
      </p>
    </div>

    {#if successMsg}
      <div class="rounded-md bg-green-50 border border-green-200 px-4 py-3 text-sm text-green-700 dark:bg-green-950 dark:border-green-800 dark:text-green-300">
        {successMsg}
      </div>
    {/if}

    {#if error}
      <div class="rounded-md bg-destructive/10 border border-destructive/20 px-4 py-3 text-sm text-destructive">
        {error}
      </div>
    {/if}

    {#if mode === 'login'}
      <form onsubmit={handleLogin} class="space-y-4">
        <div class="space-y-2">
          <Label for="email">{t('login.email')}</Label>
          <Input
            id="email"
            type="email"
            bind:value={email}
            placeholder={t('login.emailPlaceholder')}
            required
            autocomplete="email"
          />
        </div>

        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <Label for="password">{t('login.password')}</Label>
            <button
              type="button"
              class="text-xs text-primary hover:underline"
              onclick={goToForgot}
            >
              {t('login.forgotPassword')}
            </button>
          </div>
          <Input
            id="password"
            type="password"
            bind:value={password}
            placeholder={t('login.passwordPlaceholder')}
            required
            autocomplete="current-password"
          />
        </div>

        <Button type="submit" class="w-full" disabled={loading || !email || !password}>
          {loading ? t('login.submitting') : t('login.submit')}
        </Button>
      </form>

    {:else if mode === 'forgot'}
      <form onsubmit={handleSendResetCode} class="space-y-4">
        <div class="space-y-2">
          <Label for="reset-email">{t('forgot.email')}</Label>
          <Input
            id="reset-email"
            type="email"
            bind:value={email}
            placeholder={t('forgot.emailPlaceholder')}
            required
            autocomplete="email"
          />
        </div>

        <Button type="submit" class="w-full" disabled={loading || !email || cooldown > 0}>
          {#if loading}
            {t('forgot.sendingCode')}
          {:else if cooldown > 0}
            {t('forgot.pleaseWait', { cooldown })}
          {:else}
            {t('forgot.sendCode')}
          {/if}
        </Button>
      </form>

      <button
        type="button"
        class="text-sm text-muted-foreground hover:underline w-full text-center"
        onclick={backToLogin}
      >
        {t('forgot.backToLogin')}
      </button>

    {:else}
      <form onsubmit={handleResetPassword} class="space-y-4">
        <p class="text-sm text-muted-foreground">
          {@html t('forgot.codeSent', { email })}
        </p>

        <div class="space-y-2">
          <Label for="reset-code">{t('forgot.verificationCode')}</Label>
          <Input
            id="reset-code"
            type="text"
            bind:value={resetCode}
            placeholder={t('forgot.codePlaceholder')}
            required
            maxlength={6}
            inputmode="numeric"
            autocomplete="one-time-code"
          />
        </div>

        <div class="space-y-2">
          <Label for="new-password">{t('forgot.newPassword')}</Label>
          <Input
            id="new-password"
            type="password"
            bind:value={newPassword}
            placeholder={t('forgot.newPasswordPlaceholder')}
            required
            autocomplete="new-password"
          />
          {#if passwordTooShort}
            <p class="text-xs text-destructive">{t('forgot.passwordTooShort')}</p>
          {/if}
        </div>

        <div class="space-y-2">
          <Label for="confirm-new-password">{t('forgot.confirmPassword')}</Label>
          <Input
            id="confirm-new-password"
            type="password"
            bind:value={confirmNewPassword}
            placeholder={t('forgot.confirmPlaceholder')}
            required
            autocomplete="new-password"
          />
          {#if passwordMismatch}
            <p class="text-xs text-destructive">{t('forgot.passwordMismatch')}</p>
          {/if}
        </div>

        <Button
          type="submit"
          class="w-full"
          disabled={loading || resetCode.length !== 6 || newPassword.length < 6 || newPassword !== confirmNewPassword}
        >
          {loading ? t('forgot.submitting') : t('forgot.submit')}
        </Button>

        <div class="flex items-center justify-between">
          <button
            type="button"
            class="text-sm text-muted-foreground hover:underline"
            onclick={backToLogin}
          >
            {t('forgot.backToLogin')}
          </button>
          <button
            type="button"
            class="text-sm {cooldown > 0 ? 'text-muted-foreground cursor-not-allowed' : 'text-primary hover:underline'}"
            disabled={cooldown > 0 || loading}
            onclick={handleResendCode}
          >
            {cooldown > 0 ? t('forgot.resendCooldown', { cooldown }) : t('forgot.resend')}
          </button>
        </div>
      </form>
    {/if}

    <p class="text-center text-sm text-muted-foreground">
      {t('login.noAccount')}
      <a href="/register" class="text-primary hover:underline font-medium">{t('login.register')}</a>
    </p>
  </div>
</div>
