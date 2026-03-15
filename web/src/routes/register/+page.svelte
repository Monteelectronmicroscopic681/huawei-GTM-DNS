<script lang="ts">
  import { goto } from '$app/navigation';
  import { api } from '$lib/api';
  import { auth } from '$lib/auth.svelte';
  import { t } from '$lib/i18n.svelte';
  import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';

  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let code = $state('');
  let error = $state('');
  let loading = $state(false);
  let step = $state<'form' | 'verify'>('form');
  let cooldown = $state(0);
  let cooldownTimer: ReturnType<typeof setInterval> | null = null;

  const passwordTooShort = $derived(password.length > 0 && password.length < 6);
  const passwordMismatch = $derived(confirmPassword.length > 0 && password !== confirmPassword);
  const canSubmit = $derived(
    email.length > 0 &&
    password.length >= 6 &&
    password === confirmPassword &&
    !loading &&
    cooldown === 0
  );

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

  async function handleSendCode(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    if (password !== confirmPassword) {
      error = t('register.passwordMismatch');
      return;
    }
    loading = true;
    try {
      await api.auth.registerRequest(email, password);
      startCooldown();
      step = 'verify';
    } catch (err) {
      error = err instanceof Error ? err.message : t('register.sendFailed');
    } finally {
      loading = false;
    }
  }

  async function handleResend() {
    if (cooldown > 0 || loading) return;
    error = '';
    loading = true;
    try {
      await api.auth.registerRequest(email, password);
      startCooldown();
      code = '';
      error = '';
    } catch (err) {
      error = err instanceof Error ? err.message : t('register.resendFailed');
    } finally {
      loading = false;
    }
  }

  async function handleVerify(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    loading = true;
    try {
      const res = await api.auth.registerVerify(email, code);
      auth.login(res.token);
      goto('/');
    } catch (err) {
      error = err instanceof Error ? err.message : t('register.verifyFailed');
    } finally {
      loading = false;
    }
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
        {step === 'form' ? t('register.subtitle.form') : t('register.subtitle.verify')}
      </p>
    </div>

    {#if error}
      <div class="rounded-md bg-destructive/10 border border-destructive/20 px-4 py-3 text-sm text-destructive">
        {error}
      </div>
    {/if}

    {#if step === 'form'}
      <form onsubmit={handleSendCode} class="space-y-4">
        <div class="space-y-2">
          <Label for="email">{t('register.email')}</Label>
          <Input
            id="email"
            type="email"
            bind:value={email}
            placeholder={t('register.emailPlaceholder')}
            required
            autocomplete="email"
          />
        </div>

        <div class="space-y-2">
          <Label for="password">{t('register.password')}</Label>
          <Input
            id="password"
            type="password"
            bind:value={password}
            placeholder={t('register.passwordPlaceholder')}
            required
            autocomplete="new-password"
          />
          {#if passwordTooShort}
            <p class="text-xs text-destructive">{t('register.passwordTooShort')}</p>
          {/if}
        </div>

        <div class="space-y-2">
          <Label for="confirm-password">{t('register.confirmPassword')}</Label>
          <Input
            id="confirm-password"
            type="password"
            bind:value={confirmPassword}
            placeholder={t('register.confirmPlaceholder')}
            required
            autocomplete="new-password"
          />
          {#if passwordMismatch}
            <p class="text-xs text-destructive">{t('register.passwordMismatch')}</p>
          {/if}
        </div>

        <Button type="submit" class="w-full" disabled={!canSubmit}>
          {#if loading}
            {t('register.sendingCode')}
          {:else if cooldown > 0}
            {t('register.pleaseWait', { cooldown })}
          {:else}
            {t('register.sendCode')}
          {/if}
        </Button>
      </form>
    {:else}
      <form onsubmit={handleVerify} class="space-y-4">
        <p class="text-sm text-muted-foreground">
          {@html t('register.codeSent', { email })}
        </p>

        <div class="space-y-2">
          <Label for="code">{t('register.verificationCode')}</Label>
          <Input
            id="code"
            type="text"
            bind:value={code}
            placeholder={t('register.codePlaceholder')}
            required
            maxlength={6}
            inputmode="numeric"
            autocomplete="one-time-code"
          />
        </div>

        <Button type="submit" class="w-full" disabled={loading || code.length !== 6}>
          {loading ? t('register.verifying') : t('register.verify')}
        </Button>

        <div class="flex items-center justify-between">
          <button
            type="button"
            class="text-sm text-muted-foreground hover:underline"
            onclick={() => { step = 'form'; error = ''; code = ''; }}
          >
            {t('common.back')}
          </button>
          <button
            type="button"
            class="text-sm {cooldown > 0 ? 'text-muted-foreground cursor-not-allowed' : 'text-primary hover:underline'}"
            disabled={cooldown > 0 || loading}
            onclick={handleResend}
          >
            {cooldown > 0 ? t('register.resendCooldown', { cooldown }) : t('register.resend')}
          </button>
        </div>
      </form>
    {/if}

    <p class="text-center text-sm text-muted-foreground">
      {t('register.hasAccount')}
      <a href="/login" class="text-primary hover:underline font-medium">{t('register.signIn')}</a>
    </p>
  </div>
</div>
