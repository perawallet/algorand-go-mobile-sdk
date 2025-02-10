// Top-level build file where you can add configuration options common to all sub-projects/modules.
plugins {
    id("org.jetbrains.kotlin.android") version "2.0.21" apply false
    id("com.android.library") version "8.2.2" apply false
    `java-library`
    kotlin("jvm") version "2.0.21"
}


allprojects {
    tasks.withType<JavaCompile> {
        options.encoding = "UTF-8"
    }
}
