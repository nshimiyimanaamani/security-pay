<template functional>
  <section class="feedback">
    <header class="feedback-header">
      <div class="feedback-name">{{ props.feedback.display_name }}</div>
      <p class="feedback-date">
        {{ props.feedback.creator }}&nbsp;&nbsp;|&nbsp;&nbsp;{{
          props.feedback.update_at | date
        }}
      </p>
    </header>

    <article class="feedback-content">
      <div class="feedback-title">{{ props.feedback.title }}</div>
      {{ props.feedback.body }}
    </article>
  </section>
</template>
<style lang="scss" scoped>
.feedback {
  width: 350px;
  border: none;
  border-radius: 9px;
  padding: 0.75rem 1.5rem 1.25rem;
  background: #e5e8e9;
}

.feedback-header {
  display: flex;
  flex-direction: column;
  margin-bottom: 1rem;
  border-bottom: 2px dashed #aaa;
}

.feedback-header .feedback-name {
  margin: 0;
  font-size: 1.2rem;
  line-height: 2;
  font-weight: bold;
}
.feedback-header .feedback-date {
  font-size: 0.9rem;
}

.feedback-content {
  font-size: 0.95rem;
}
.feedback-content .feedback-title {
  line-height: 1.5;
  margin: 0;
  font-weight: bold;
  font-size: 1.1rem;
}
</style>