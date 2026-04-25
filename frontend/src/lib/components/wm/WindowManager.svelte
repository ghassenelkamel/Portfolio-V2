<script lang="ts">
  import { onMount } from "svelte";
  import { wm, wmInit } from "$lib/wm/store";
  import WMNode from "./WMNode.svelte";

  const wmState = $derived($wm);

  onMount(() => {
    wmInit();
  });
</script>

<div class="wmroot">
  <WMNode node={wmState.root} />
</div>

<style>
  .wmroot {
    width: 100%;
    height: 100%;
    padding: 10px;
    box-sizing: border-box;
  }
</style>
