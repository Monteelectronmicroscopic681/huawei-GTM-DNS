<script lang="ts">
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import { buttonVariants } from '$lib/components/ui/button';
  import { t } from '$lib/i18n.svelte';
  import { confirmState, resolveConfirm } from '$lib/confirm.svelte';
</script>

<AlertDialog.Root
  open={confirmState.open}
  onOpenChange={(v) => {
    confirmState.open = v;
    if (!v) resolveConfirm(false);
  }}
>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>{confirmState.title || t('confirm.title')}</AlertDialog.Title>
      {#if confirmState.description}
        <AlertDialog.Description>{confirmState.description}</AlertDialog.Description>
      {/if}
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>{confirmState.cancelText || t('common.cancel')}</AlertDialog.Cancel>
      <AlertDialog.Action
        class={buttonVariants({ variant: confirmState.variant })}
        onclick={() => {
          confirmState.open = false;
          resolveConfirm(true);
        }}
      >
        {confirmState.confirmText || t('common.confirm')}
      </AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
