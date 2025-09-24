<script lang="ts">
  import { cn } from "$lib/utils";
  import { onMount, tick } from "svelte";

  let { 
    containerRef = $bindable(),
    class:className = "",
    fromRef = $bindable(),
    toRef = $bindable(),
    curvature = 0,
    reverse = false, // Include the reverse prop
    pathColor = "gray",
    pathWidth = 2,
    pathOpacity = 0.2,
    startXOffset = 0,
    startYOffset = 0,
    endXOffset = 0,
    endYOffset = 0,
  }= $props();
  

  let id = crypto.randomUUID().slice(0, 8);
  let pathD = $state("");
  let svgDimensions = { width: 0, height: 0 };

  let updatePath = () => {
    if (!containerRef || !fromRef || !toRef) {
      return;
    }
    let containerRect = containerRef?.getBoundingClientRect();
    let rectA = fromRef?.getBoundingClientRect();
    let rectB = toRef?.getBoundingClientRect();

    let svgWidth = containerRect.width;
    let svgHeight = containerRect.height;
    svgDimensions.width = svgWidth;
    svgDimensions.height = svgHeight;

    let startX =
      rectA.left - containerRect.left + rectA.width / 2 + startXOffset;
    let startY =
      rectA.top - containerRect.top + rectA.height / 2 + startYOffset;
    let endX = rectB.left - containerRect.left + rectB.width / 2 + endXOffset;
    let endY = rectB.top - containerRect.top + rectB.height / 2 + endYOffset;

    let controlY = startY - curvature;
    let d = `M ${startX},${startY} Q ${
      (startX + endX) / 2
    },${controlY} ${endX},${endY}`;
    pathD = d;
  };
  onMount(async () => {
    await tick().then(() => {
      updatePath();
      const resizeObserver = new ResizeObserver((entries) => {
        // For all entries, recalculate the path
        for (let entry of entries) {
          updatePath();
        }
      });

      // Observe the container element
      if (containerRef) {
        resizeObserver.observe(containerRef);
      }
    });
  });
</script>

<svg
  fill="none"
  width={svgDimensions.width}
  height={svgDimensions.height}
  xmlns="http://www.w3.org/2000/svg"
  class={cn(
    "pointer-events-none absolute left-0 top-0 transform-gpu stroke-2 animate-pulse",
    className,
  )}
  viewBox={`0 0 ${svgDimensions.width} ${svgDimensions.height}`}
>
  <path
    d={pathD}
    stroke={pathColor}
    stroke-width={pathWidth}
    stroke-opacity={pathOpacity}
    stroke-linecap="round"
  />
  <path
    d={pathD}
    stroke-width={pathWidth}
    stroke={`url(#${id})`}
    stroke-opacity="1"
    stroke-linecap="round"
  />
  <defs>
    <linearGradient {id} gradientUnits="userSpaceOnUse" class="transform-gpu">
      <stop class="[stop-color:var(--color-primary)]" stop-opacity="0"></stop>
      <stop class="[stop-color:var(--color-primary)]"></stop>
      <stop offset="32.5%" class="[stop-color:var(--color-primary)]"></stop>
      <stop offset="100%" class="[stop-color:var(--color-primary)]" stop-opacity="0"
      ></stop>
    </linearGradient>
  </defs>
</svg>
