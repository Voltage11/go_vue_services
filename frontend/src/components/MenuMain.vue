<script setup>
const menuItems = [
  {
    title: 'Главная',
    href: '/',
    children: [],
  },
  {
    title: 'Расписание',
    href: '#',
    children: [],
  },
  {
    title: 'Клиенты',
    href: '#',
    children: [],
  },
  {
    title: 'Услуги',
    href: '#',
    children: [],
  },

  {
    title: 'Персонал',
    href: '#',
    children: [],
  },
]
</script>

<template>
  <div v-for="menuItem in menuItems" :key="menuItem">
    <div class="menu-item-block">
      <div class="menu-item-block__top">
        <router-link :to="menuItem.href" class="link-group">{{ menuItem.title.toUpperCase() }}</router-link>

      </div>
      <div  v-if="menuItem.children.length > 0" class="menu-item-block__buttom">
        <div v-for="subItem in menuItem.children" :key="subItem">
          <div>
            <router-link :to="subItem.href" class="link-item">{{ subItem.title }}</router-link>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.menu-item-block{
  border-right: solid 1px var(--link-color-light-hover);
  padding-right: 15px;
}

.menu-item-block__top {
  align-items: center;
  display: flex;
  justify-content: center;
}

.menu-item-block__buttom {
  display: flex;
  gap: 15px;
}

.link-group, .link-item{
  font-size: 13px;
  text-decoration: none;
  color: var(--link-color-light-hover)
}

.link-group{
  font-weight: 600;
}

.link-group:hover, link-item:hover{
  color: var(--link-color-light)
}

</style>
