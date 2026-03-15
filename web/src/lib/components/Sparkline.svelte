<script lang="ts">
  let { data = [], width = 80, height = 20 }: { data: (0|1)[], width?: number, height?: number } = $props();

  const points = $derived(
    data.map((v, i) => {
      const x = (i / Math.max(data.length - 1, 1)) * width;
      const y = v === 1 ? 4 : height - 4;
      return `${x},${y}`;
    }).join(' ')
  );

  const color = $derived(data.length > 0 && data[data.length - 1] === 1 ? '#22c55e' : '#ef4444');
</script>

<svg {width} {height} class="inline-block">
  {#if data.length > 1}
    <polyline {points} fill="none" stroke={color} stroke-width="1.5" />
  {/if}
</svg>
