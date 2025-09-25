<script lang="ts"></script>

<div>
  <p>
    Eleakxir's search engine is designed to be both fast and flexible, letting
    you find what you need in multiple ways.
  </p>

  <h3 class="h3 mt-4 mb-2">Search Modes</h3>

  <p>
    <span class="text-primary font-semibold">All:</span> This is the default mode.
    It searches for your query across a set of standard columns like email, username,
    and phone number. This is the fastest and most efficient way to find a specific
    user or account.
  </p>

  <p>
    <span class="text-primary font-semibold">Specific column:</span>
    This mode lets you choose a specific column to search within, such as email or
    username. It's useful when you know exactly where the data you're looking for
    is stored.
  </p>

  <p>
    <span class="text-primary font-semibold">Full Text:</span> This mode combines
    all available columns into a single, large text field and searches within it.
    It's great for finding data that might be in an unexpected column, but it's way
    slower.
  </p>

  <h3 class="h3 mt-4 mb-2">Query Matching</h3>
  <p>
    <span class="text-primary font-semibold">Standard Search:</span> By default,
    Eleakxir uses a "fuzzy" search. This means it will find results where your search
    terms are part of a larger string. For example, searching for 1234 would find
    john.doe@1234.com. Multiple terms are combined with <b>AND</b>, so all terms
    must be present, while multiple columns are combined with <b>OR</b>, so a match
    in any column counts.
  </p>

  <p>
    <span class="text-primary font-semibold">Exact Match:</span> When you enable
    "Exact Match," the search will only return results where the data in a column
    is an exact match for your search term. This is useful for finding specific,
    unique values.
  </p>

  <p>
    <span class="text-primary font-semibold">Starts With or Ends With:</span>
    You can refine your search by using <code>^</code> at the beginning of a term
    to match only values that <b>start with</b> that term, or <code>$</code> at the
    end of a term to match only values that <b>end with</b> that term. For example,
    <code>^admin</code> will find entries starting with "admin", while
    <code>user$</code> will find entries ending with "user".
  </p>
</div>
